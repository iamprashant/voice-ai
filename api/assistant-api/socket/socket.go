// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package assistant_socket

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rapidaai/api/assistant-api/config"
	internal_adapter "github.com/rapidaai/api/assistant-api/internal/adapters"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony"
	internal_assistant_service "github.com/rapidaai/api/assistant-api/internal/services/assistant"
	web_client "github.com/rapidaai/pkg/clients/web"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/storages"
	storage_files "github.com/rapidaai/pkg/storages/file-storage"
	"github.com/rapidaai/pkg/utils"
)

// AudioSocketEngine manages Asterisk AudioSocket TCP connections.
// It delegates call-context resolution to InboundDispatcher.ResolveCallSessionByContext
// to avoid duplicating entity-loading logic.
type audioSocketEngine struct {
	logger   commons.Logger
	cfg      *config.AssistantConfig
	listener net.Listener
	mu       sync.RWMutex

	postgres   connectors.PostgresConnector
	redis      connectors.RedisConnector
	opensearch connectors.OpenSearchConnector
	storage    storages.Storage

	inboundDispatcher *internal_telephony.InboundDispatcher
}

// NewAudioSocketEngine creates a new AudioSocket engine.
// Internally builds an InboundDispatcher for call-context resolution (Postgres
// lookup + atomic claiming). The engine only manages TCP + streamer.
func NewAudioSocketEngine(config *config.AssistantConfig, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector,
) *audioSocketEngine {
	store := callcontext.NewStore(postgres, logger)
	vaultClient := web_client.NewVaultClientGRPC(&config.AppConfig, logger, redis)
	fileStorage := storage_files.NewStorage(config.AssetStoreConfig, logger)
	assistantService := internal_assistant_service.NewAssistantService(config, logger, postgres, opensearch)
	conversationService := internal_assistant_service.NewAssistantConversationService(logger, postgres, fileStorage)

	dispatcher := internal_telephony.NewInboundDispatcher(internal_telephony.TelephonyDispatcherDeps{
		Cfg:                 config,
		Logger:              logger,
		Store:               store,
		VaultClient:         vaultClient,
		AssistantService:    assistantService,
		ConversationService: conversationService,
	})

	return &audioSocketEngine{
		cfg:               config,
		logger:            logger,
		postgres:          postgres,
		redis:             redis,
		opensearch:        opensearch,
		storage:           fileStorage,
		inboundDispatcher: dispatcher,
	}
}

// Start begins the AudioSocket TCP listener.
func (m *audioSocketEngine) Connect(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.AudioSocketConfig.Host, m.cfg.AudioSocketConfig.Port)
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
func (m *audioSocketEngine) Disconnect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listener == nil {
		return nil
	}
	_ = m.listener.Close()
	m.listener = nil
	return nil
}

func (m *audioSocketEngine) acceptLoop(ctx context.Context) {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			m.logger.Warnw("AudioSocket accept error", "error", err)
			continue
		}

		go m.handleConnection(ctx, conn)
	}
}

func (m *audioSocketEngine) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Step 1: Read the UUID frame from AudioSocket protocol.
	// Asterisk sends a FrameTypeUUID (0x01) as the first frame with a 16-byte UUID payload.
	// This UUID is the contextId that was passed via the dialplan (e.g. AudioSocket(contextId, host:port)).
	contextID, err := m.readContextID(reader)
	if err != nil {
		m.logger.Warnw("AudioSocket failed to read UUID frame", "error", err)
		return
	}

	m.logger.Infof("AudioSocket connection received with contextId=%s", contextID)

	// Step 2: Resolve call context — delegates to InboundDispatcher which handles
	// Postgres lookup, atomic claiming, and parallel entity loading.
	cc, vaultCred, err := m.inboundDispatcher.ResolveCallSessionByContext(connCtx, contextID)
	if err != nil {
		m.logger.Warnw("AudioSocket session resolution failed", "contextId", contextID, "error", err)
		return
	}

	// Step 3: Create AudioSocket streamer and start talking.
	// Pass the contextID as the initial UUID so the streamer sends ConversationInitialization
	// on the first Recv() call — the UUID frame was already consumed by readContextID above.
	streamer, err := internal_telephony.Telephony(internal_telephony.Asterisk).NewStreamer(
		m.logger, cc, vaultCred, internal_telephony.StreamerOption{
			AudioSocketConn:   conn,
			AudioSocketReader: reader,
			AudioSocketWriter: writer,
		},
	)
	if err != nil {
		m.logger.Warnw("AudioSocket streamer create failed", "contextId", contextID, "error", err)
		return
	}

	talker, err := internal_adapter.GetTalker(
		utils.PhoneCall,
		connCtx,
		m.cfg,
		m.logger,
		m.postgres,
		m.opensearch,
		m.redis,
		m.storage,
		streamer,
	)
	if err != nil {
		m.logger.Warnw("AudioSocket talker create failed", "contextId", contextID, "error", err)
		return
	}

	if err := talker.Talk(connCtx, cc.ToAuth()); err != nil {
		m.logger.Warnw("AudioSocket talker exited", "contextId", contextID, "error", err)
	}

	// Mark call context as completed now that the call has ended
	m.inboundDispatcher.CompleteCallSession(connCtx, contextID)
}

// readContextID reads the initial UUID frame from the AudioSocket connection.
// Asterisk sends a FrameTypeUUID (0x01) frame with 16-byte UUID payload.
// Frame format: 1-byte type + 2-byte big-endian length + payload.
// We parse the UUID and return it as a string (the contextId).
func (m *audioSocketEngine) readContextID(reader *bufio.Reader) (string, error) {
	const frameTypeUUID byte = 0x01

	// Read frame type
	frameType, err := reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("failed to read frame type: %w", err)
	}

	// Read 2-byte big-endian length
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(reader, lenBuf); err != nil {
		return "", fmt.Errorf("failed to read frame length: %w", err)
	}
	payloadLen := int(binary.BigEndian.Uint16(lenBuf))

	// Read payload
	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(reader, payload); err != nil {
			return "", fmt.Errorf("failed to read frame payload: %w", err)
		}
	}

	if frameType != frameTypeUUID {
		return "", fmt.Errorf("expected UUID frame (0x01), got frame type 0x%02x", frameType)
	}

	if len(payload) != 16 {
		return "", fmt.Errorf("invalid UUID payload length: %d (expected 16)", len(payload))
	}

	// Convert 16-byte UUID to standard string format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	h := hex.EncodeToString(payload)
	uuid := h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]

	return uuid, nil
}
