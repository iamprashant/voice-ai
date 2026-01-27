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
)

type websocketExecutor struct {
	logger       commons.Logger
	conn         *websocket.Conn
	history      []*protos.Message
	mu           sync.RWMutex
	writeMu      sync.Mutex
	done         chan struct{}
	closed       bool
	requestTimes sync.Map
}

// NewWebsocketAssistantExecutor creates a new WebSocket-based assistant executor.
func NewWebsocketAssistantExecutor(logger commons.Logger) internal_agent_executor.AssistantExecutor {
	return &websocketExecutor{
		logger:  logger,
		history: make([]*protos.Message, 0),
		done:    make(chan struct{}),
	}
}

// Name returns the executor name identifier.
func (e *websocketExecutor) Name() string {
	return "websocket"
}

// Initialize establishes the WebSocket connection and starts the listener.
func (e *websocketExecutor) Initialize(
	ctx context.Context,
	comm internal_type.Communication,
	config *protos.AssistantConversationConfiguration,
) error {
	start := time.Now()
	_, span, _ := comm.Tracer().StartSpan(ctx, utils.AssistantAgentConnectStage,
		internal_adapter_telemetry.KV{K: "executor", V: internal_adapter_telemetry.StringValue(e.Name())})
	defer span.EndSpan(ctx, utils.AssistantAgentConnectStage)

	provider := comm.Assistant().AssistantProviderWebsocket
	if provider == nil {
		return fmt.Errorf("websocket provider is nil")
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
	if err := e.send(Request{
		Type:      TypeConfiguration,
		Timestamp: time.Now().UnixMilli(),
		Data: ConfigurationData{
			AssistantID:    comm.Assistant().Id,
			ConversationID: comm.Conversation().Id,
		},
	}); err != nil {
		return fmt.Errorf("failed to send configuration: %w", err)
	}

	e.logger.Benchmark("WebsocketExecutor.Initialize", time.Since(start))
	return nil
}

// connect establishes the WebSocket connection.
func (e *websocketExecutor) connect(ctx context.Context, provider *internal_assistant_entity.AssistantProviderWebsocket) error {
	headers := http.Header{}
	for k, v := range provider.Headers {
		headers.Set(k, v)
	}

	wsURL, err := url.Parse(provider.Url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	query := wsURL.Query()
	for k, v := range provider.Parameters {
		query.Set(k, v)
	}
	wsURL.RawQuery = query.Encode()

	dialer := websocket.Dialer{HandshakeTimeout: 30 * time.Second}
	conn, _, err := dialer.DialContext(ctx, wsURL.String(), headers)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	conn.SetReadLimit(10 * 1024 * 1024)
	e.conn = conn
	return nil
}

// send writes a message to the WebSocket.
func (e *websocketExecutor) send(msg Request) error {
	e.writeMu.Lock()
	defer e.writeMu.Unlock()

	if e.conn == nil {
		return fmt.Errorf("not connected")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return e.conn.WriteMessage(websocket.TextMessage, data)
}

// listen reads messages from WebSocket until context is cancelled or connection closes.
func (e *websocketExecutor) listen(ctx context.Context, comm internal_type.Communication) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-e.done:
			return nil
		default:
		}

		if e.conn == nil {
			return fmt.Errorf("not connected")
		}

		// Allow periodic context checks
		e.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

		_, data, err := e.conn.ReadMessage()
		if err != nil {
			if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
				continue
			}
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return fmt.Errorf("server_closed")
			}
			return err
		}

		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			e.logger.Errorf("Invalid response: %v", err)
			continue
		}

		e.handleResponse(ctx, &resp, comm)
	}
}

// handleResponse processes a single response from the server.
// Server sends responses sequentially - one at a time per request.
func (e *websocketExecutor) handleResponse(ctx context.Context, resp *Response, comm internal_type.Communication) {
	if resp.Error != nil {
		e.logger.Errorf("Server error: %d - %s", resp.Error.Code, resp.Error.Message)
		return
	}

	convID := fmt.Sprintf("%d", comm.Conversation().Id)
	getID := func(id string) string {
		if id != "" {
			return id
		}
		return convID
	}

	switch resp.Type {
	case TypeError:
		var d ErrorData
		json.Unmarshal(resp.Data, &d)
		e.logger.Errorf("Error: %d - %s", d.Code, d.Message)

	case TypeStream:
		// Streaming chunk - forward to assistant
		var d StreamData
		json.Unmarshal(resp.Data, &d)
		comm.OnPacket(ctx, internal_type.LLMStreamPacket{ContextID: getID(d.ID), Text: d.Content})

	case TypeComplete:
		// Response complete - store in history and send metrics
		var d CompleteData
		json.Unmarshal(resp.Data, &d)
		id := getID(d.ID)

		if d.Content != "" {
			msg := types.NewMessage("assistant", &types.Content{
				ContentType:   commons.TEXT_CONTENT.String(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Content:       []byte(d.Content),
			})
			e.mu.Lock()
			e.history = append(e.history, msg.ToProto())
			e.mu.Unlock()
			comm.OnPacket(ctx, internal_type.LLMMessagePacket{ContextID: id, Message: msg})
		}

		// Send metrics
		var metrics []*types.Metric
		if t, ok := e.requestTimes.LoadAndDelete(id); ok {
			metrics = append(metrics, types.NewTimeTakenMetric(time.Since(t.(time.Time))))
		}
		for _, m := range d.Metrics {
			metrics = append(metrics, &types.Metric{Name: m.Name, Value: fmt.Sprintf("%f", m.Value), Description: m.Unit})
		}
		if len(metrics) > 0 {
			comm.OnPacket(ctx, internal_type.MetricPacket{ContextID: id, Metrics: metrics})
		}

	case TypeToolCall:
		// Server requests action (disconnect, hold, etc)
		var d ToolCallData
		json.Unmarshal(resp.Data, &d)
		action := e.mapToolAction(d.Name)
		comm.OnPacket(ctx, internal_type.LLMToolPacket{
			ContextID: getID(d.ID),
			Name:      d.Name,
			Action:    action,
			Result:    d.Params,
		})

	case TypeInterruption:
		// User interrupted the response
		var d InterruptionData
		json.Unmarshal(resp.Data, &d)
		source := internal_type.InterruptionSourceWord
		if d.Source == "vad" {
			source = internal_type.InterruptionSourceVad
		}
		comm.OnPacket(ctx, internal_type.InterruptionPacket{ContextID: getID(d.ID), Source: source})

	case TypeClose:
		// Server closed the session
		var d CloseData
		json.Unmarshal(resp.Data, &d)
		e.closed = true
		comm.OnPacket(ctx, internal_type.ClosePacket{Reason: d.Reason})

	case TypePing:
		e.send(Request{Type: TypePong, Timestamp: time.Now().UnixMilli()})
	}
}

// mapToolAction maps tool names from websocket to conversation actions.
func (e *websocketExecutor) mapToolAction(name string) protos.AssistantConversationAction_ActionType {
	switch name {
	case "disconnect", "end_conversation", "hangup":
		return protos.AssistantConversationAction_END_CONVERSATION
	case "hold", "put_on_hold":
		return protos.AssistantConversationAction_PUT_ON_HOLD
	default:
		return protos.AssistantConversationAction_ACTION_UNSPECIFIED
	}
}

// Execute sends a packet to the WebSocket server.
func (e *websocketExecutor) Execute(ctx context.Context, comm internal_type.Communication, packet internal_type.Packet) error {
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

		return e.send(Request{
			Type:      TypeUserMessage,
			Timestamp: time.Now().UnixMilli(),
			Data:      UserMessageData{ID: id, Content: p.Text},
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

// Close terminates the WebSocket connection.
func (e *websocketExecutor) Close(ctx context.Context, comm internal_type.Communication) error {
	select {
	case <-e.done:
	default:
		close(e.done)
	}

	if e.conn != nil && !e.closed {
		e.writeMu.Lock()
		e.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		e.writeMu.Unlock()
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
