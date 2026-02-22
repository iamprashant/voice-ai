// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_agentkit

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	internal_agent_executor "github.com/rapidaai/api/assistant-api/internal/agent/executor"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ internal_agent_executor.AssistantExecutor = (*agentkitExecutor)(nil)

type agentkitExecutor struct {
	logger     commons.Logger
	connection *grpc.ClientConn
	talker     grpc.BidiStreamingClient[protos.TalkInput, protos.TalkOutput]
	mu         sync.RWMutex
}

// NewAgentKitAssistantExecutor creates a new AgentKit-based assistant executor.
func NewAgentKitAssistantExecutor(logger commons.Logger) internal_agent_executor.AssistantExecutor {
	return &agentkitExecutor{
		logger: logger,
	}
}

// Name returns the executor name identifier.
func (e *agentkitExecutor) Name() string {
	return "agentkit"
}

// Initialize establishes the gRPC connection and starts the listener.
func (e *agentkitExecutor) Initialize(ctx context.Context, comm internal_type.Communication, cfg *protos.ConversationInitialization) error {
	_, span, _ := comm.Tracer().StartSpan(ctx, utils.AssistantAgentConnectStage, internal_adapter_telemetry.KV{K: "executor", V: internal_adapter_telemetry.StringValue(e.Name())})
	defer span.EndSpan(ctx, utils.AssistantAgentConnectStage)

	provider := comm.Assistant().AssistantProviderAgentkit
	if provider == nil {
		return fmt.Errorf("agentkit provider is not enabled")
	}

	// Connect
	if err := e.connect(ctx, provider); err != nil {
		return err
	}

	// Start listener - stops on context cancel or server close
	utils.Go(ctx, func() {
		if err := e.listen(ctx, comm.OnPacket); err != nil && ctx.Err() == nil {
			comm.OnPacket(ctx, internal_type.DirectivePacket{Directive: protos.ConversationDirective_END_CONVERSATION, Arguments: map[string]interface{}{"reason": err.Error()}})
		}
	})

	// Send initial configuration
	if err := e.sendConfiguration(provider.AssistantId, provider.Id, comm.Conversation().Id, cfg); err != nil {
		return fmt.Errorf("failed to send configuration: %w", err)
	}
	return nil
}

// connect establishes the gRPC connection.
func (e *agentkitExecutor) connect(ctx context.Context, provider *internal_assistant_entity.AssistantProviderAgentkit) error {
	opts := []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64), grpc.MaxCallSendMsgSize(math.MaxInt64))}
	// credentials and tls
	if provider.Certificate != "" {
		creds, err := e.buildTLSCredentials(provider.Certificate)
		if err != nil {
			return fmt.Errorf("TLS credentials failed: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// create connection with provider url
	conn, err := grpc.NewClient(provider.Url, opts...)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	// create client and stream
	client := protos.NewAgentKitClient(conn)
	e.connection = conn

	// Build metadata from provider.Metadata (headers to pass to server)
	streamCtx := ctx
	if len(provider.Metadata) > 0 {
		md := metadata.New(map[string]string(provider.Metadata))
		streamCtx = metadata.NewOutgoingContext(ctx, md)
	}

	talker, err := client.Talk(streamCtx)
	if err != nil {
		return fmt.Errorf("stream start failed: %w", err)
	}
	e.talker = talker
	return nil
}

// buildTLSCredentials creates TLS credentials from a PEM certificate.
// If certPEM is "insecure" or "skip-verify", it skips certificate verification (dev only).
func (e *agentkitExecutor) buildTLSCredentials(certPEM string) (credentials.TransportCredentials, error) {
	// Allow skipping verification for development
	if certPEM == "insecure" || certPEM == "skip-verify" {
		e.logger.Warnf("Using insecure TLS (skipping certificate verification) - DO NOT USE IN PRODUCTION")
		return credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}), nil
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(certPEM)) {
		e.logger.Errorf("Failed to parse certificate PEM (length=%d, starts=%q)", len(certPEM), certPEM[:min(50, len(certPEM))])
		return nil, fmt.Errorf("invalid certificate: failed to parse PEM")
	}
	return credentials.NewTLS(&tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}), nil
}

// send writes a message to the gRPC stream.
func (e *agentkitExecutor) send(req *protos.TalkInput) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.talker == nil {
		return fmt.Errorf("not connected")
	}
	return e.talker.Send(req)
}

// sendConfiguration sends the initial configuration.
func (e *agentkitExecutor) sendConfiguration(assistantId uint64, assistantProviderID uint64, ConversationID uint64, cfg *protos.ConversationInitialization) error {
	return e.send(&protos.TalkInput{
		Request: &protos.TalkInput_Initialization{
			Initialization: &protos.ConversationInitialization{
				AssistantConversationId: ConversationID,
				Assistant: &protos.AssistantDefinition{
					AssistantId: assistantId,
					Version:     utils.GetVersionString(assistantProviderID),
				},
				Args:         cfg.GetArgs(),
				Metadata:     cfg.GetMetadata(),
				Options:      cfg.GetOptions(),
				StreamMode:   cfg.GetStreamMode(),
				UserIdentity: cfg.GetUserIdentity(),
				Time:         timestamppb.New(time.Now()),
			},
		},
	})
}

// listen reads messages until context is cancelled or connection closes.
func (e *agentkitExecutor) listen(ctx context.Context, onPacket func(ctx context.Context, packet ...internal_type.Packet) error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		resp, err := e.talker.Recv()
		if err != nil {
			e.logger.Debugf("Listener received error: %v", err)
			code := status.Code(err)
			switch {
			case errors.Is(err, io.EOF):
				// Server gracefully closed
				onPacket(ctx, internal_type.DirectivePacket{Directive: protos.ConversationDirective_END_CONVERSATION, Arguments: map[string]interface{}{"reason": "server closed connection"}})
			case code == codes.Canceled:
				// RPC canceled (client or server)
				onPacket(ctx, internal_type.DirectivePacket{Directive: protos.ConversationDirective_END_CONVERSATION, Arguments: map[string]interface{}{"reason": "connection canceled"}})
			case code == codes.Unavailable:
				// Server went down
				onPacket(ctx, internal_type.DirectivePacket{Directive: protos.ConversationDirective_END_CONVERSATION, Arguments: map[string]interface{}{"reason": "server unavailable"}})
			default:
				// Other errors
				onPacket(ctx, internal_type.DirectivePacket{Directive: protos.ConversationDirective_END_CONVERSATION, Arguments: map[string]interface{}{"reason": err.Error()}})
			}
			return nil
		}
		e.handleResponse(ctx, resp, onPacket)
	}
}

// handleResponse processes a single response from the server.
func (e *agentkitExecutor) handleResponse(ctx context.Context, resp *protos.TalkOutput, onPacket func(ctx context.Context, packet ...internal_type.Packet) error) {
	switch data := resp.GetData().(type) {
	case *protos.TalkOutput_Interruption:
		onPacket(ctx, internal_type.InterruptionPacket{ContextID: data.Interruption.Id, Source: internal_type.InterruptionSourceWord})

	case *protos.TalkOutput_Assistant:
		switch msg := data.Assistant.GetMessage().(type) {
		case *protos.ConversationAssistantMessage_Text:
			if data.Assistant.GetCompleted() {
				onPacket(ctx, internal_type.LLMResponseDonePacket{
					ContextID: data.Assistant.GetId(),
					Text:      msg.Text,
				})
				return
			}
			onPacket(ctx, internal_type.LLMResponseDeltaPacket{ContextID: data.Assistant.GetId(), Text: msg.Text})
		case *protos.ConversationAssistantMessage_Audio:
			e.logger.Debugf("Received audio message (not implemented)")
		}

	case *protos.TalkOutput_Directive:
		args, _ := utils.AnyMapToInterfaceMap(data.Directive.GetArgs())
		onPacket(ctx, internal_type.DirectivePacket{ContextID: data.Directive.GetId(), Directive: data.Directive.GetType(), Arguments: args})
	}
}

// Execute sends a packet to the AgentKit server.
func (e *agentkitExecutor) Execute(ctx context.Context, comm internal_type.Communication, packet internal_type.Packet) error {
	_, span, _ := comm.Tracer().StartSpan(ctx, utils.AssistantAgentTextGenerationStage, internal_adapter_telemetry.MessageKV(packet.ContextId()))
	defer span.EndSpan(ctx, utils.AssistantAgentTextGenerationStage)
	switch p := packet.(type) {
	case internal_type.UserTextPacket:
		return e.send(&protos.TalkInput{
			Request: &protos.TalkInput_Message{
				Message: &protos.ConversationUserMessage{
					Message: &protos.ConversationUserMessage_Text{
						Text: p.Text,
					},
					Id:        packet.ContextId(),
					Completed: true,
					Time:      timestamppb.Now(),
				},
			},
		})
	case internal_type.StaticPacket:
		return nil

	default:
		return fmt.Errorf("unsupported packet: %T", packet)
	}
}

// Close terminates the gRPC connection.
func (e *agentkitExecutor) Close(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.talker != nil {
		e.talker.CloseSend()
		e.talker = nil
	}
	if e.connection != nil {
		e.connection.Close()
		e.connection = nil
	}
	return nil
}
