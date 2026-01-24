// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	internal_agent_executor "github.com/rapidaai/api/assistant-api/internal/agent/executor"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"golang.org/x/sync/errgroup"
)

// =============================================================================
// WebSocket Message Types - Similar to AgentKit gRPC pattern
// =============================================================================

// WSMessageType defines the type of message and what data structure to expect
type WSMessageType string

const (
	// Request types (client -> server)
	WSTypeConfiguration WSMessageType = "configuration" // Data: WSConfigurationData
	WSTypeUserMessage   WSMessageType = "user_message"  // Data: WSUserMessageData

	// Response types (server -> client)
	WSTypeAssistantMessage WSMessageType = "assistant_message" // Data: WSAssistantMessageData
	WSTypeStream           WSMessageType = "stream"            // Data: WSStreamData
	WSTypeInterruption     WSMessageType = "interruption"      // Data: WSInterruptionData
	WSTypeError            WSMessageType = "error"             // Data: WSErrorData

	// Control types (bidirectional)
	WSTypePing WSMessageType = "ping" // Data: nil
	WSTypePong WSMessageType = "pong" // Data: nil
)

// =============================================================================
// Request/Response Envelope
// =============================================================================

// WSRequest represents an outgoing WebSocket message with typed data
type WSRequest struct {
	Type      WSMessageType `json:"type"`
	Timestamp int64         `json:"timestamp"`
	Data      interface{}   `json:"data,omitempty"`
}

// WSResponse represents an incoming WebSocket message with typed data
type WSResponse struct {
	Type    WSMessageType   `json:"type"`
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *WSErrorData    `json:"error,omitempty"`
}

// =============================================================================
// Data Structures for each message type
// =============================================================================

// WSConfigurationData contains initial connection configuration
// Used with: WSTypeConfiguration
type WSConfigurationData struct {
	AssistantID         uint64                 `json:"assistant_id"`
	ConversationID      uint64                 `json:"conversation_id"`
	AssistantDefinition *WSAssistantDefinition `json:"assistant,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Args                map[string]interface{} `json:"args,omitempty"`
	Options             map[string]interface{} `json:"options,omitempty"`
}

// WSAssistantDefinition contains assistant metadata
type WSAssistantDefinition struct {
	AssistantID uint64 `json:"assistant_id"`
	Name        string `json:"name,omitempty"`
}

// WSUserMessageData contains user message content
// Used with: WSTypeUserMessage
type WSUserMessageData struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Timestamp int64  `json:"timestamp"`
}

// WSAssistantMessageData contains assistant response content
// Used with: WSTypeAssistantMessage
type WSAssistantMessageData struct {
	ID        string                     `json:"id"`
	Message   *WSAssistantMessageContent `json:"message"`
	Completed bool                       `json:"completed"`
	Metrics   []*WSMetric                `json:"metrics,omitempty"`
}

// WSAssistantMessageContent represents the message content (text or audio)
type WSAssistantMessageContent struct {
	Type    string `json:"type"` // "text" or "audio"
	Content string `json:"content,omitempty"`
	Audio   []byte `json:"audio,omitempty"`
}

// WSStreamData contains streaming text chunk
// Used with: WSTypeStream
type WSStreamData struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Index   int    `json:"index,omitempty"`
}

// WSInterruptionData contains interruption signal
// Used with: WSTypeInterruption
type WSInterruptionData struct {
	ID      string  `json:"id,omitempty"`
	Source  string  `json:"source,omitempty"` // "word", "vad"
	StartAt float64 `json:"start_at,omitempty"`
	EndAt   float64 `json:"end_at,omitempty"`
}

// WSErrorData contains error information
// Used with: WSTypeError or in WSResponse.Error
type WSErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WSMetric contains metric information
type WSMetric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

type websocketExecutor struct {
	logger     commons.Logger
	connection *websocket.Conn
	history    []*protos.Message
	mu         sync.RWMutex
	writeMu    sync.Mutex // Separate mutex for write operations
	done       chan struct{}
}

// NewWebsocketAssistantExecutor creates a new WebSocket-based assistant executor
// that communicates with an external HTTPS/WebSocket service for LLM operations.
func NewWebsocketAssistantExecutor(logger commons.Logger) internal_agent_executor.AssistantExecutor {
	return &websocketExecutor{
		logger:  logger,
		history: make([]*protos.Message, 0),
		done:    make(chan struct{}),
	}
}

// Name returns the executor name identifier.
func (executor *websocketExecutor) Name() string {
	return "websocket"
}

// Initialize establishes the WebSocket connection and starts the message listener.
func (executor *websocketExecutor) Initialize(
	ctx context.Context,
	communication internal_type.Communication,
	config *protos.AssistantConversationConfiguration,
) error {
	start := time.Now()
	ctx, span, _ := communication.Tracer().StartSpan(
		ctx,
		utils.AssistantAgentConnectStage,
		internal_adapter_telemetry.KV{K: "executor", V: internal_adapter_telemetry.StringValue(executor.Name())},
	)
	defer span.EndSpan(ctx, utils.AssistantAgentConnectStage)

	providerDefinition := communication.Assistant().AssistantProviderWebsocket
	if providerDefinition == nil {
		return fmt.Errorf("websocket provider definition is nil")
	}

	g, gCtx := errgroup.WithContext(ctx)

	// Goroutine to establish WebSocket connection
	g.Go(func() error {
		return executor.establishConnection(gCtx, providerDefinition)
	})

	// Goroutine to fetch conversation history
	g.Go(func() error {
		executor.mu.Lock()
		defer executor.mu.Unlock()
		executor.history = append(executor.history, communication.GetConversationLogs()...)
		return nil
	})

	if err := g.Wait(); err != nil {
		executor.logger.Errorf("Error during initialization of websocket: %v", err)
		return err
	}

	// Start the response listener in background
	utils.Go(ctx, func() {
		if err := executor.responseListener(ctx, communication); err != nil {
			executor.logger.Errorf("Error in WebSocket response listener: %v", err)
		}
	})

	// Send initial configuration
	if err := executor.sendConfiguration(communication); err != nil {
		return fmt.Errorf("failed to send configuration: %w", err)
	}

	executor.logger.Benchmark("WebsocketExecutor.Initialize", time.Since(start))
	return nil
}

// establishConnection creates the WebSocket connection with proper headers and parameters.
func (executor *websocketExecutor) establishConnection(
	ctx context.Context,
	provider *internal_assistant_entity.AssistantProviderWebsocket,
) error {
	// Prepare HTTP headers
	headers := http.Header{}
	if provider.Headers != nil {
		for key, value := range provider.Headers {
			headers.Set(key, value)
		}
	}

	// Parse and modify WebSocket URL
	wsURL, err := url.Parse(provider.Url)
	if err != nil {
		return fmt.Errorf("failed to parse websocket URL: %w", err)
	}

	// Add query parameters
	query := wsURL.Query()
	if provider.Parameters != nil {
		for key, value := range provider.Parameters {
			query.Set(key, value)
		}
		wsURL.RawQuery = query.Encode()
	}

	// Configure dialer with timeout
	dialer := websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, wsURL.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	// Configure connection settings
	conn.SetReadLimit(10 * 1024 * 1024) // 10MB max message size
	conn.SetPongHandler(func(appData string) error {
		executor.logger.Debugf("Received pong from WebSocket server")
		return nil
	})
	executor.connection = conn
	return nil
}

// sendConfiguration sends the initial configuration to the WebSocket service.
func (executor *websocketExecutor) sendConfiguration(communication internal_type.Communication) error {
	config := WSRequest{
		Type:      WSTypeConfiguration,
		Timestamp: time.Now().UnixMilli(),
		Data: WSConfigurationData{
			AssistantID:    communication.Assistant().Id,
			ConversationID: communication.Conversation().Id,
			AssistantDefinition: &WSAssistantDefinition{
				AssistantID: communication.Assistant().Id,
			},
			Metadata: map[string]interface{}{
				"history_length": len(executor.history),
			},
		},
	}

	return executor.sendMessage(config)
}

// sendMessage safely sends a message over the WebSocket connection.
func (executor *websocketExecutor) sendMessage(msg WSRequest) error {
	executor.writeMu.Lock()
	defer executor.writeMu.Unlock()

	if executor.connection == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	executor.logger.Debugf("Sending WebSocket message: type=%s", msg.Type)
	if err := executor.connection.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// responseListener listens for responses from the WebSocket service and processes them.
func (executor *websocketExecutor) responseListener(ctx context.Context, communication internal_type.Communication) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-executor.done:
			return nil
		default:
		}

		if executor.connection == nil {
			return fmt.Errorf("websocket connection is nil")
		}

		_, message, err := executor.connection.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				executor.logger.Debugf("WebSocket connection closed normally")
				return nil
			}
			return fmt.Errorf("websocket read error: %w", err)
		}

		var resp WSResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			executor.logger.Errorf("Failed to unmarshal WebSocket response: %v", err)
			continue
		}

		executor.logger.Debugf("Received WebSocket message: type=%s, success=%v", resp.Type, resp.Success)
		if err := executor.processResponse(ctx, &resp, communication); err != nil {
			executor.logger.Errorf("Error processing WebSocket response: %v", err)
		}
	}
}

// processResponse handles individual responses from the WebSocket service.
// The response type determines what data structure is in the Data field.
func (executor *websocketExecutor) processResponse(
	ctx context.Context,
	resp *WSResponse,
	communication internal_type.Communication,
) error {
	// Handle error response
	if resp.Error != nil {
		executor.logger.Errorf("WebSocket error response: code=%d, message=%s, details=%s",
			resp.Error.Code, resp.Error.Message, resp.Error.Details)
		return nil
	}

	switch resp.Type {
	case WSTypeError:
		// Parse error data
		var errorData WSErrorData
		if err := json.Unmarshal(resp.Data, &errorData); err != nil {
			executor.logger.Errorf("Failed to parse error data: %v", err)
			return nil
		}
		executor.logger.Errorf("WebSocket error: code=%d, message=%s", errorData.Code, errorData.Message)
		return nil

	case WSTypeStream:
		// Parse stream data
		var streamData WSStreamData
		if err := json.Unmarshal(resp.Data, &streamData); err != nil {
			executor.logger.Errorf("Failed to parse stream data: %v", err)
			return nil
		}
		contextID := streamData.ID
		if contextID == "" {
			contextID = fmt.Sprintf("%d", communication.Conversation().Id)
		}
		communication.OnPacket(ctx, internal_type.LLMStreamPacket{
			ContextID: contextID,
			Text:      streamData.Content,
		})

	case WSTypeAssistantMessage:
		// Parse assistant message data
		var msgData WSAssistantMessageData
		if err := json.Unmarshal(resp.Data, &msgData); err != nil {
			executor.logger.Errorf("Failed to parse assistant message data: %v", err)
			return nil
		}
		contextID := msgData.ID
		if contextID == "" {
			contextID = fmt.Sprintf("%d", communication.Conversation().Id)
		}

		// Handle text message
		if msgData.Message != nil && msgData.Message.Type == "text" && msgData.Message.Content != "" {
			message := types.NewMessage("assistant", &types.Content{
				ContentType:   commons.TEXT_CONTENT.String(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Content:       []byte(msgData.Message.Content),
			})

			executor.mu.Lock()
			executor.history = append(executor.history, message.ToProto())
			executor.mu.Unlock()

			communication.OnPacket(ctx, internal_type.LLMMessagePacket{
				ContextID: contextID,
				Message:   message,
			})
		}

		// Handle metrics if present
		if len(msgData.Metrics) > 0 {
			metrics := make([]*types.Metric, 0, len(msgData.Metrics))
			for _, m := range msgData.Metrics {
				metrics = append(metrics, &types.Metric{
					Name:        m.Name,
					Value:       fmt.Sprintf("%f", m.Value),
					Description: m.Unit,
				})
			}
			communication.OnPacket(ctx, internal_type.MetricPacket{
				ContextID: contextID,
				Metrics:   metrics,
			})
		}

	case WSTypeInterruption:
		// Parse interruption data
		var interruptData WSInterruptionData
		if err := json.Unmarshal(resp.Data, &interruptData); err != nil {
			executor.logger.Errorf("Failed to parse interruption data: %v", err)
			return nil
		}
		contextID := interruptData.ID
		if contextID == "" {
			contextID = fmt.Sprintf("%d", communication.Conversation().Id)
		}
		source := internal_type.InterruptionSourceWord
		if interruptData.Source == "vad" {
			source = internal_type.InterruptionSourceVad
		}
		communication.OnPacket(ctx, internal_type.InterruptionPacket{
			ContextID: contextID,
			Source:    source,
			StartAt:   interruptData.StartAt,
			EndAt:     interruptData.EndAt,
		})

	case WSTypePing:
		// Respond with pong
		executor.sendMessage(WSRequest{
			Type:      WSTypePong,
			Timestamp: time.Now().UnixMilli(),
		})

	case WSTypePong:
		executor.logger.Debugf("Received pong message")
	}

	return nil
}

// Execute processes incoming packets and sends them to the WebSocket service.
func (executor *websocketExecutor) Execute(
	ctx context.Context,
	communication internal_type.Communication,
	packet internal_type.Packet,
) error {
	ctx, span, _ := communication.Tracer().StartSpan(
		ctx,
		utils.AssistantAgentTextGenerationStage,
		internal_adapter_telemetry.MessageKV(packet.ContextId()),
	)
	defer span.EndSpan(ctx, utils.AssistantAgentTextGenerationStage)

	switch p := packet.(type) {
	case internal_type.UserTextPacket:
		return executor.handleUserTextPacket(ctx, p, communication)
	case internal_type.StaticPacket:
		return executor.handleStaticPacket(p)
	default:
		return fmt.Errorf("unsupported packet type: %T", packet)
	}
}

// handleUserTextPacket sends user text messages to the WebSocket service.
func (executor *websocketExecutor) handleUserTextPacket(
	ctx context.Context,
	packet internal_type.UserTextPacket,
	communication internal_type.Communication,
) error {
	// Record user message in history
	userMessage := types.NewMessage("user", &types.Content{
		ContentType:   commons.TEXT_CONTENT.String(),
		ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
		Content:       []byte(packet.Text),
	})

	executor.mu.Lock()
	executor.history = append(executor.history, userMessage.ToProto())
	executor.mu.Unlock()

	// Send message over WebSocket with typed data
	msg := WSRequest{
		Type:      WSTypeUserMessage,
		Timestamp: time.Now().UnixMilli(),
		Data: WSUserMessageData{
			ID:        packet.ContextId(),
			Content:   packet.Text,
			Completed: true,
			Timestamp: time.Now().UnixMilli(),
		},
	}

	return executor.sendMessage(msg)
}

// handleStaticPacket appends static assistant responses to history.
func (executor *websocketExecutor) handleStaticPacket(packet internal_type.StaticPacket) error {
	executor.mu.Lock()
	defer executor.mu.Unlock()

	executor.history = append(executor.history, &protos.Message{
		Role: "assistant",
		Contents: []*protos.Content{
			{
				ContentType:   commons.TEXT_CONTENT.String(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Content:       []byte(packet.Text),
			},
		},
	})
	return nil
}

// Close terminates the WebSocket connection and cleans up resources.
func (executor *websocketExecutor) Close(ctx context.Context, communication internal_type.Communication) error {
	executor.logger.Debugf("Closing WebSocket executor")

	// Signal done to stop the listener
	close(executor.done)

	if executor.connection != nil {
		// Send close message
		executor.writeMu.Lock()
		err := executor.connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		executor.writeMu.Unlock()

		if err != nil {
			executor.logger.Errorf("Error sending close message: %v", err)
		}

		if err := executor.connection.Close(); err != nil {
			executor.logger.Errorf("Error closing WebSocket connection: %v", err)
		}
		executor.connection = nil
	}

	executor.mu.Lock()
	executor.history = make([]*protos.Message, 0)
	executor.mu.Unlock()

	// Reset done channel for potential reuse
	executor.done = make(chan struct{})

	return nil
}
