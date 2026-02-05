// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package assistant_talk_api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/rapidaai/api/assistant-api/config"
	internal_adapter "github.com/rapidaai/api/assistant-api/internal/adapters"
	internal_audiosocket "github.com/rapidaai/api/assistant-api/internal/audiosocket"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
)

// AudioSocketManager manages Asterisk AudioSocket TCP connections.
type AudioSocketManager struct {
	logger   commons.Logger
	cApi     *ConversationApi
	config   *config.AudioSocketConfig
	listener net.Listener
	mu       sync.RWMutex
}

// NewAudioSocketManager creates a new AudioSocket manager.
func NewAudioSocketManager(cApi *ConversationApi, cfg *config.AudioSocketConfig) *AudioSocketManager {
	return &AudioSocketManager{
		logger: cApi.logger,
		cApi:   cApi,
		config: cfg,
	}
}

// Start begins the AudioSocket TCP listener.
func (m *AudioSocketManager) Start(ctx context.Context) error {
	if m.config == nil || !m.config.Enabled {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("audiosocket listen failed: %w", err)
	}
	m.listener = listener

	m.logger.Info("AudioSocket server started", "addr", addr)

	go m.acceptLoop(ctx)
	return nil
}

// Close stops the AudioSocket listener.
func (m *AudioSocketManager) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listener == nil {
		return nil
	}
	_ = m.listener.Close()
	m.listener = nil
	return nil
}

func (m *AudioSocketManager) acceptLoop(ctx context.Context) {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			m.logger.Warn("AudioSocket accept error", "error", err)
			continue
		}

		go m.handleConnection(ctx, conn)
	}
}

func (m *AudioSocketManager) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// channelID, err := m.readInitialUUID(reader)
	// if err != nil {
	// 	m.logger.Warn("AudioSocket missing UUID", "error", err)
	// 	return
	// }

	auth, assistant, conversation, callerID, err := m.resolveCallContext(connCtx, "channelID")
	if err != nil {
		m.logger.Warn("AudioSocket context resolution failed", "error", err)
		return
	}

	streamer, err := internal_audiosocket.NewStreamer(m.logger, conn, reader, writer, assistant, conversation, nil)
	if err != nil {
		m.logger.Warn("AudioSocket streamer create failed", "error", err)
		return
	}

	if s, ok := streamer.(*internal_audiosocket.Streamer); ok {
		s.SetInitialUUID("channelID")
	}

	identifier := internal_adapter.Identifier(utils.PhoneCall, connCtx, auth, callerID)

	talker, err := internal_adapter.GetTalker(
		utils.PhoneCall,
		connCtx,
		m.cApi.cfg,
		m.logger,
		m.cApi.postgres,
		m.cApi.opensearch,
		m.cApi.redis,
		m.cApi.storage,
		streamer,
	)
	if err != nil {
		m.logger.Warn("AudioSocket talker create failed", "error", err)
		return
	}

	if err := talker.Talk(connCtx, auth, identifier); err != nil {
		m.logger.Warn("AudioSocket talker exited", "error", err)
	}

}

func (m *AudioSocketManager) readInitialUUID(reader *bufio.Reader) (string, error) {
	for {
		frame, err := internal_audiosocket.ReadFrame(reader)
		if err != nil {
			return "", err
		}
		switch frame.Type {
		case internal_audiosocket.FrameTypeUUID:
			return strings.TrimSpace(string(frame.Payload)), nil
		case internal_audiosocket.FrameTypeHangup:
			return "", io.EOF
		case internal_audiosocket.FrameTypeAudio:
			// Ignore audio until UUID is received
		default:
			// Ignore other frames
		}
	}
}

// resolveCallContext parses the UUID which contains apiKey:assistantId format
// Example: rpd-prj-xxx:123 or rpd-prj-xxx:123:callerID
func (m *AudioSocketManager) resolveCallContext(ctx context.Context, uuidParam string) (types.SimplePrinciple, *internal_assistant_entity.Assistant, *internal_conversation_entity.AssistantConversation, string, error) {

	auth := &types.ServiceScope{
		ProjectId:      utils.Ptr(uint64(2257831930382778368)),
		OrganizationId: utils.Ptr(uint64(2257831925018263552)),
		CurrentToken:   "3dd5c2eef53d27942bccd892750fda23ea0b92965d4699e73d8e754ab882955f",
	}

	assistant, err := m.cApi.assistantService.Get(ctx, auth, 2263072539095859200, utils.GetVersionDefinition("latest"), &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		return nil, nil, nil, "", err
	}

	conversation, err := m.cApi.assistantConversationService.CreateConversation(
		ctx,
		auth,
		internal_adapter.Identifier(utils.PhoneCall, ctx, auth, "2263072539095859200"),
		assistant.Id,
		assistant.AssistantProviderId,
		type_enums.DIRECTION_INBOUND,
		utils.PhoneCall,
	)
	if err != nil {
		return nil, nil, nil, "", err
	}

	_, _ = m.cApi.assistantConversationService.ApplyConversationMetadata(ctx, auth, assistant.Id, conversation.Id,
		[]*types.Metadata{types.NewMetadata("telephony.uuid", uuidParam)})

	return auth, assistant, conversation, "2263072539095859200", nil
}
