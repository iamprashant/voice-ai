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
	"github.com/rapidaai/pkg/types"
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

type agentkitExecutor struct {
	logger       commons.Logger
	client       protos.AgentKitClient
	talker       grpc.BidiStreamingClient[protos.AssistantMessagingRequest, protos.AssistantMessagingResponse]
	conn         *grpc.ClientConn
	history      []*protos.Message
	mu           sync.RWMutex
	done         chan struct{}
	closed       bool
	requestTimes sync.Map
}

// NewAgentKitAssistantExecutor creates a new AgentKit-based assistant executor.
func NewAgentKitAssistantExecutor(logger commons.Logger) internal_agent_executor.AssistantExecutor {
	return &agentkitExecutor{
		logger:  logger,
		history: make([]*protos.Message, 0),
		done:    make(chan struct{}),
	}
}

// Name returns the executor name identifier.
func (e *agentkitExecutor) Name() string {
	return "agentkit"
}

// Initialize establishes the gRPC connection and starts the listener.
func (e *agentkitExecutor) Initialize(ctx context.Context, comm internal_type.Communication, cfg *protos.AssistantConversationConfiguration) error {
	start := time.Now()
	_, span, _ := comm.Tracer().StartSpan(ctx, utils.AssistantAgentConnectStage,
		internal_adapter_telemetry.KV{K: "executor", V: internal_adapter_telemetry.StringValue(e.Name())})
	defer span.EndSpan(ctx, utils.AssistantAgentConnectStage)

	provider := comm.Assistant().AssistantProviderAgentkit
	if provider == nil {
		return fmt.Errorf("agentkit provider is nil")
	}

	// Connect
	if err := e.connect(ctx, provider); err != nil {
		return err
	}

	// Load history
	e.mu.Lock()
	e.history = append(e.history, comm.GetConversationLogs()...)
	e.mu.Unlock()

	// Start listener - stops on context cancel or server close
	utils.Go(ctx, func() {
		err := e.listen(ctx, comm)
		if err != nil && ctx.Err() == nil {
			e.logger.Errorf("Listener error: %v", err)
			comm.OnPacket(ctx, internal_type.ClosePacket{Reason: err.Error()})
		}
	})

	// Send configuration
	if err := e.sendConfiguration(comm, cfg); err != nil {
		return fmt.Errorf("failed to send configuration: %w", err)
	}

	e.logger.Benchmark("AgentKitExecutor.Initialize", time.Since(start))
	return nil
}

// connect establishes the gRPC connection.
func (e *agentkitExecutor) connect(ctx context.Context, provider *internal_assistant_entity.AssistantProviderAgentkit) error {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt64),
			grpc.MaxCallSendMsgSize(math.MaxInt64),
		),
	}

	// Configure TLS if certificate is provided
	if provider.Certificate != "" {
		creds, err := e.buildTLSCredentials(provider.Certificate)
		if err != nil {
			return fmt.Errorf("TLS credentials failed: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(provider.Url, opts...)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	e.conn = conn
	e.client = protos.NewAgentKitClient(conn)

	// Build metadata from provider.Metadata (headers to pass to server)
	streamCtx := ctx
	if len(provider.Metadata) > 0 {
		md := metadata.New(map[string]string(provider.Metadata))
		streamCtx = metadata.NewOutgoingContext(ctx, md)
	}

	talker, err := e.client.Talk(streamCtx)
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
func (e *agentkitExecutor) send(req *protos.AssistantMessagingRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.talker == nil {
		return fmt.Errorf("not connected")
	}
	return e.talker.Send(req)
}

// sendConfiguration sends the initial configuration.
func (e *agentkitExecutor) sendConfiguration(comm internal_type.Communication, cfg *protos.AssistantConversationConfiguration) error {
	return e.send(&protos.AssistantMessagingRequest{
		Request: &protos.AssistantMessagingRequest_Configuration{
			Configuration: &protos.AssistantConversationConfiguration{
				AssistantConversationId: comm.Conversation().Id,
				Assistant: &protos.AssistantDefinition{
					AssistantId: comm.Assistant().Id,
					Version:     utils.GetVersionString(comm.Assistant().AssistantProviderId),
				},
				Args:         cfg.GetArgs(),
				Metadata:     cfg.GetMetadata(),
				Options:      cfg.GetOptions(),
				InputConfig:  cfg.GetInputConfig(),
				OutputConfig: cfg.GetOutputConfig(),
			},
		},
	})
}

// listen reads messages until context is cancelled or connection closes.
func (e *agentkitExecutor) listen(ctx context.Context, comm internal_type.Communication) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-e.done:
			return nil
		default:
		}

		if e.talker == nil {
			return fmt.Errorf("not connected")
		}

		resp, err := e.talker.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || status.Code(err) == codes.Canceled {
				e.closed = true
				return fmt.Errorf("server_closed")
			}
			return fmt.Errorf("recv error: %w", err)
		}

		e.handleResponse(ctx, resp, comm)
	}
}

// handleResponse processes a single response from the server.
func (e *agentkitExecutor) handleResponse(ctx context.Context, resp *protos.AssistantMessagingResponse, comm internal_type.Communication) {
	if resp.GetError() != nil {
		e.logger.Errorf("Server error: code=%d, message=%s", resp.GetError().GetErrorCode(), resp.GetError().GetErrorMessage())
		return
	}

	if !resp.GetSuccess() {
		return
	}

	convID := fmt.Sprintf("%d", comm.Conversation().Id)
	getID := func(id string) string {
		if id != "" {
			return id
		}
		return convID
	}

	switch data := resp.GetData().(type) {
	case *protos.AssistantMessagingResponse_Interruption:
		// User interrupted
		comm.OnPacket(ctx, internal_type.InterruptionPacket{
			ContextID: convID,
			Source:    internal_type.InterruptionSourceWord,
		})

	case *protos.AssistantMessagingResponse_Assistant:
		e.processAssistantMessage(ctx, data.Assistant, comm, getID)

	case *protos.AssistantMessagingResponse_Action:
		// Server requests an action (tool call)
		action := data.Action
		comm.OnPacket(ctx, internal_type.LLMToolPacket{
			ContextID: convID,
			Name:      action.GetName(),
			Action:    action.GetAction(),
			Result:    nil,
		})
	}
}

// processAssistantMessage handles assistant messages.
func (e *agentkitExecutor) processAssistantMessage(
	ctx context.Context,
	assistant *protos.AssistantConversationAssistantMessage,
	comm internal_type.Communication,
	getID func(string) string,
) {
	if assistant == nil {
		return
	}

	id := getID(assistant.GetId())

	switch msg := assistant.GetMessage().(type) {
	case *protos.AssistantConversationAssistantMessage_Text:
		content := msg.Text.GetContent()

		// Send streaming chunk
		comm.OnPacket(ctx, internal_type.LLMStreamPacket{
			ContextID: id,
			Text:      content,
		})

		// If completed, store in history and send full message
		if assistant.GetCompleted() {
			message := types.NewMessage("assistant", &types.Content{
				ContentType:   commons.TEXT_CONTENT.String(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Content:       []byte(content),
			})

			e.mu.Lock()
			e.history = append(e.history, message.ToProto())
			e.mu.Unlock()

			comm.OnPacket(ctx, internal_type.LLMMessagePacket{
				ContextID: id,
				Message:   message,
			})

			// Send metrics
			var metrics []*types.Metric
			if t, ok := e.requestTimes.LoadAndDelete(id); ok {
				metrics = append(metrics, types.NewTimeTakenMetric(time.Since(t.(time.Time))))
			}
			if len(metrics) > 0 {
				comm.OnPacket(ctx, internal_type.MetricPacket{ContextID: id, Metrics: metrics})
			}
		}

	case *protos.AssistantConversationAssistantMessage_Audio:
		e.logger.Debugf("Received audio message (not implemented)")
	}
}

// mapToolAction maps tool names to conversation actions.
func (e *agentkitExecutor) mapToolAction(name string) protos.AssistantConversationAction_ActionType {
	switch name {
	case "disconnect", "end_conversation", "hangup":
		return protos.AssistantConversationAction_END_CONVERSATION
	case "hold", "put_on_hold":
		return protos.AssistantConversationAction_PUT_ON_HOLD
	default:
		return protos.AssistantConversationAction_ACTION_UNSPECIFIED
	}
}

// Execute sends a packet to the AgentKit server.
func (e *agentkitExecutor) Execute(ctx context.Context, comm internal_type.Communication, packet internal_type.Packet) error {
	_, span, _ := comm.Tracer().StartSpan(ctx, utils.AssistantAgentTextGenerationStage,
		internal_adapter_telemetry.MessageKV(packet.ContextId()))
	defer span.EndSpan(ctx, utils.AssistantAgentTextGenerationStage)

	switch p := packet.(type) {
	case internal_type.UserTextPacket:
		if e.closed {
			return fmt.Errorf("connection closed")
		}

		id := p.ContextId()
		e.requestTimes.Store(id, time.Now())

		// Store in history
		msg := types.NewMessage("user", &types.Content{
			ContentType:   commons.TEXT_CONTENT.String(),
			ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
			Content:       []byte(p.Text),
		})
		e.mu.Lock()
		e.history = append(e.history, msg.ToProto())
		e.mu.Unlock()

		return e.send(&protos.AssistantMessagingRequest{
			Request: &protos.AssistantMessagingRequest_Message{
				Message: &protos.AssistantConversationUserMessage{
					Message: &protos.AssistantConversationUserMessage_Text{
						Text: &protos.AssistantConversationMessageTextContent{
							Content: p.Text,
						},
					},
					Id:        id,
					Completed: true,
					Time:      timestamppb.Now(),
				},
			},
		})

	case internal_type.StaticPacket:
		e.mu.Lock()
		e.history = append(e.history, &protos.Message{
			Role:     "assistant",
			Contents: []*protos.Content{{ContentType: commons.TEXT_CONTENT.String(), ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(), Content: []byte(p.Text)}},
		})
		e.mu.Unlock()
		return nil

	default:
		return fmt.Errorf("unsupported packet: %T", packet)
	}
}

// Close terminates the gRPC connection.
func (e *agentkitExecutor) Close(ctx context.Context, comm internal_type.Communication) error {
	select {
	case <-e.done:
	default:
		close(e.done)
	}

	if e.talker != nil {
		e.talker.CloseSend()
		e.talker = nil
	}

	if e.conn != nil {
		e.conn.Close()
		e.conn = nil
	}

	e.mu.Lock()
	e.history = nil
	e.closed = false
	e.mu.Unlock()

	e.done = make(chan struct{})
	return nil
}
