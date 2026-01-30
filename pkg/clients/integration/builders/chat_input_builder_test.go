// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package integration_client_builders

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rapidaai/pkg/commons"
	gorm_types "github.com/rapidaai/pkg/models/gorm/types"
	"github.com/rapidaai/protos"
)

// mockLogger implements commons.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Level() zapcore.Level                                           { return zapcore.DebugLevel }
func (m *mockLogger) Debug(args ...interface{})                                      {}
func (m *mockLogger) Debugf(template string, args ...interface{})                    {}
func (m *mockLogger) Info(args ...interface{})                                       {}
func (m *mockLogger) Infof(template string, args ...interface{})                     {}
func (m *mockLogger) Warn(args ...interface{})                                       {}
func (m *mockLogger) Warnf(template string, args ...interface{})                     {}
func (m *mockLogger) Error(args ...interface{})                                      {}
func (m *mockLogger) Errorf(template string, args ...interface{})                    {}
func (m *mockLogger) DPanic(args ...interface{})                                     {}
func (m *mockLogger) DPanicf(template string, args ...interface{})                   {}
func (m *mockLogger) Panic(args ...interface{})                                      {}
func (m *mockLogger) Panicf(template string, args ...interface{})                    {}
func (m *mockLogger) Fatal(args ...interface{})                                      {}
func (m *mockLogger) Fatalf(template string, args ...interface{})                    {}
func (m *mockLogger) Benchmark(functionName string, duration time.Duration)          {}
func (m *mockLogger) Tracef(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) Sync() error                                                    { return nil }

var _ commons.Logger = (*mockLogger)(nil)

func newTestLogger() commons.Logger {
	return &mockLogger{}
}

func TestNewChatInputBuilder(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	assert.NotNil(t, builder, "NewChatInputBuilder should return a non-nil builder")
	assert.Implements(t, (*InputChatBuilder)(nil), builder, "builder should implement InputChatBuilder")
}

func TestChatInputBuilder_Credential(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	tests := []struct {
		name     string
		id       uint64
		value    map[string]interface{}
		expected uint64
	}{
		{
			name:     "basic credential",
			id:       12345,
			value:    map[string]interface{}{"api_key": "test-key"},
			expected: 12345,
		},
		{
			name:     "zero id credential",
			id:       0,
			value:    map[string]interface{}{},
			expected: 0,
		},
		{
			name:     "large id credential",
			id:       18446744073709551615, // max uint64
			value:    map[string]interface{}{"key": "value", "nested": map[string]interface{}{"a": "b"}},
			expected: 18446744073709551615,
		},
		{
			name:     "credential with multiple fields",
			id:       999,
			value:    map[string]interface{}{"api_key": "key123", "secret": "secret456", "region": "us-east-1"},
			expected: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structVal, err := structpb.NewStruct(tt.value)
			require.NoError(t, err, "failed to create struct")

			cred := builder.Credential(tt.id, structVal)

			assert.NotNil(t, cred, "credential should not be nil")
			assert.Equal(t, tt.expected, cred.Id, "credential id should match")
			assert.NotNil(t, cred.Value, "credential value should not be nil")
		})
	}
}

func TestChatInputBuilder_Credential_NilValue(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	cred := builder.Credential(123, nil)

	assert.NotNil(t, cred, "credential should not be nil")
	assert.Equal(t, uint64(123), cred.Id, "credential id should match")
	assert.Nil(t, cred.Value, "credential value should be nil")
}

func TestChatInputBuilder_WithinMessage(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger).(*inputChatBuilder)

	tests := []struct {
		name         string
		role         string
		prompt       string
		expectedRole string
		checkFunc    func(*testing.T, *protos.Message)
	}{
		{
			name:         "user message",
			role:         "user",
			prompt:       "Hello, world!",
			expectedRole: "user",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				user := msg.GetUser()
				require.NotNil(t, user, "user message should not be nil")
				assert.Equal(t, "Hello, world!", user.GetContent(), "user content should match")
			},
		},
		{
			name:         "system message",
			role:         "system",
			prompt:       "You are a helpful assistant.",
			expectedRole: "system",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				system := msg.GetSystem()
				require.NotNil(t, system, "system message should not be nil")
				assert.Equal(t, "You are a helpful assistant.", system.GetContent(), "system content should match")
			},
		},
		{
			name:         "assistant message",
			role:         "assistant",
			prompt:       "I'm here to help!",
			expectedRole: "assistant",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				assistant := msg.GetAssistant()
				require.NotNil(t, assistant, "assistant message should not be nil")
				require.Len(t, assistant.GetContents(), 1, "assistant should have one content")
				assert.Equal(t, "I'm here to help!", assistant.GetContents()[0], "assistant content should match")
			},
		},
		{
			name:         "unknown role defaults to assistant",
			role:         "unknown",
			prompt:       "Default content",
			expectedRole: "unknown",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				assistant := msg.GetAssistant()
				require.NotNil(t, assistant, "unknown role should default to assistant message")
				require.Len(t, assistant.GetContents(), 1, "assistant should have one content")
				assert.Equal(t, "Default content", assistant.GetContents()[0], "content should match")
			},
		},
		{
			name:         "empty role defaults to assistant",
			role:         "",
			prompt:       "Empty role content",
			expectedRole: "",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				assistant := msg.GetAssistant()
				require.NotNil(t, assistant, "empty role should default to assistant message")
			},
		},
		{
			name:         "user message with empty content",
			role:         "user",
			prompt:       "",
			expectedRole: "user",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				user := msg.GetUser()
				require.NotNil(t, user, "user message should not be nil")
				assert.Equal(t, "", user.GetContent(), "user content should be empty")
			},
		},
		{
			name:         "user message with special characters",
			role:         "user",
			prompt:       "Hello! @#$%^&*() 你好 مرحبا",
			expectedRole: "user",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				user := msg.GetUser()
				require.NotNil(t, user, "user message should not be nil")
				assert.Equal(t, "Hello! @#$%^&*() 你好 مرحبا", user.GetContent(), "special characters should be preserved")
			},
		},
		{
			name:         "user message with newlines",
			role:         "user",
			prompt:       "Line 1\nLine 2\nLine 3",
			expectedRole: "user",
			checkFunc: func(t *testing.T, msg *protos.Message) {
				user := msg.GetUser()
				require.NotNil(t, user, "user message should not be nil")
				assert.Equal(t, "Line 1\nLine 2\nLine 3", user.GetContent(), "newlines should be preserved")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := builder.WithinMessage(tt.role, tt.prompt)

			assert.NotNil(t, msg, "message should not be nil")
			assert.Equal(t, tt.expectedRole, msg.GetRole(), "role should match")
			tt.checkFunc(t, msg)
		})
	}
}

func TestChatInputBuilder_Chat(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger).(*inputChatBuilder)

	t.Run("basic chat request", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		modelOpts := make(map[string]*anypb.Any)
		tempVal, _ := structpb.NewValue(0.7)
		anyTemp, _ := anypb.New(tempVal)
		modelOpts["temperature"] = anyTemp

		conversations := []*protos.Message{
			builder.WithinMessage("user", "Hello"),
		}

		request := builder.Chat(
			"req-123",
			credential,
			modelOpts,
			nil,
			nil,
			conversations...,
		)

		assert.NotNil(t, request, "request should not be nil")
		assert.Equal(t, "req-123", request.RequestId, "request id should match")
		assert.NotNil(t, request.Credential, "credential should not be nil")
		assert.Len(t, request.Conversations, 1, "should have one conversation")
		assert.NotNil(t, request.ModelParameters, "model parameters should not be nil")
		assert.Nil(t, request.ToolDefinitions, "tool definitions should be nil")
	})

	t.Run("chat request with tools", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		tools := []*protos.FunctionDefinition{
			{
				Name:        "get_weather",
				Description: "Get the current weather",
			},
			{
				Name:        "search",
				Description: "Search the web",
			},
		}

		request := builder.Chat(
			"req-456",
			credential,
			nil,
			tools,
			nil,
			builder.WithinMessage("user", "What's the weather?"),
		)

		assert.NotNil(t, request, "request should not be nil")
		require.Len(t, request.ToolDefinitions, 2, "should have two tool definitions")
		assert.Equal(t, "function", request.ToolDefinitions[0].Type, "tool type should be function")
		assert.Equal(t, "get_weather", request.ToolDefinitions[0].FunctionDefinition.Name, "function name should match")
		assert.Equal(t, "search", request.ToolDefinitions[1].FunctionDefinition.Name, "function name should match")
	})

	t.Run("chat request with additional data", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		additionalData := map[string]string{
			"trace_id": "trace-123",
			"user_id":  "user-456",
		}

		request := builder.Chat(
			"req-789",
			credential,
			nil,
			nil,
			additionalData,
		)

		assert.NotNil(t, request, "request should not be nil")
		assert.Equal(t, "trace-123", request.AdditionalData["trace_id"], "trace_id should match")
		assert.Equal(t, "user-456", request.AdditionalData["user_id"], "user_id should match")
	})

	t.Run("chat request with multiple conversations", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		conversations := []*protos.Message{
			builder.WithinMessage("system", "You are a helpful assistant"),
			builder.WithinMessage("user", "Hello"),
			builder.WithinMessage("assistant", "Hi there!"),
			builder.WithinMessage("user", "How are you?"),
		}

		request := builder.Chat(
			"req-multi",
			credential,
			nil,
			nil,
			nil,
			conversations...,
		)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Conversations, 4, "should have four conversations")
		assert.Equal(t, "system", request.Conversations[0].GetRole(), "first message should be system")
		assert.Equal(t, "user", request.Conversations[1].GetRole(), "second message should be user")
		assert.Equal(t, "assistant", request.Conversations[2].GetRole(), "third message should be assistant")
		assert.Equal(t, "user", request.Conversations[3].GetRole(), "fourth message should be user")
	})

	t.Run("chat request with no conversations", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		request := builder.Chat(
			"req-empty",
			credential,
			nil,
			nil,
			nil,
		)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Conversations, 0, "should have no conversations")
	})
}

func TestChatInputBuilder_Message(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	t.Run("single template", func(t *testing.T) {
		templates := []*gorm_types.PromptTemplate{
			{Role: "user", Content: "Hello, {{name}}!"},
		}
		args := map[string]interface{}{"name": "World"}

		messages := builder.Message(templates, args)

		require.Len(t, messages, 1, "should have one message")
		assert.Equal(t, "user", messages[0].GetRole(), "role should be user")
		user := messages[0].GetUser()
		require.NotNil(t, user, "user message should not be nil")
		assert.Equal(t, "Hello, World!", user.GetContent(), "content should be templated")
	})

	t.Run("multiple templates", func(t *testing.T) {
		templates := []*gorm_types.PromptTemplate{
			{Role: "system", Content: "You are {{assistant_type}}."},
			{Role: "user", Content: "My name is {{name}}."},
		}
		args := map[string]interface{}{
			"assistant_type": "a helpful assistant",
			"name":           "Alice",
		}

		messages := builder.Message(templates, args)

		require.Len(t, messages, 2, "should have two messages")
		assert.Equal(t, "system", messages[0].GetRole(), "first role should be system")
		assert.Equal(t, "user", messages[1].GetRole(), "second role should be user")
	})

	t.Run("empty templates", func(t *testing.T) {
		templates := []*gorm_types.PromptTemplate{}
		args := map[string]interface{}{}

		messages := builder.Message(templates, args)

		assert.Len(t, messages, 0, "should have no messages")
	})

	t.Run("template with missing argument", func(t *testing.T) {
		templates := []*gorm_types.PromptTemplate{
			{Role: "user", Content: "Hello, {{name}}!"},
		}
		args := map[string]interface{}{} // missing "name"

		messages := builder.Message(templates, args)

		require.Len(t, messages, 1, "should have one message")
		// The template parser should handle missing arguments gracefully
	})

	t.Run("nil arguments", func(t *testing.T) {
		templates := []*gorm_types.PromptTemplate{
			{Role: "user", Content: "Hello, World!"},
		}

		messages := builder.Message(templates, nil)

		require.Len(t, messages, 1, "should have one message")
	})
}

func TestChatInputBuilder_Arguments(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	t.Run("merge with existing variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "name", DefaultValue: "default_name"},
			{Name: "age", DefaultValue: "25"},
		}

		nameVal, _ := structpb.NewValue("Alice")
		nameAny, _ := anypb.New(nameVal)

		arguments := map[string]*anypb.Any{
			"name": nameAny,
		}

		result := builder.Arguments(variables, arguments)

		// The anypb value gets converted to a map with protobuf metadata
		assert.NotNil(t, result["name"], "name should be overridden")
		assert.Equal(t, "25", result["age"], "age should use default")
	})

	t.Run("empty variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}

		nameVal, _ := structpb.NewValue("Alice")
		nameAny, _ := anypb.New(nameVal)

		arguments := map[string]*anypb.Any{
			"name": nameAny,
		}

		result := builder.Arguments(variables, arguments)

		// Check the key exists - actual value structure depends on protobuf conversion
		assert.NotNil(t, result["name"], "name should be set from arguments")
	})

	t.Run("nil arguments", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "name", DefaultValue: "default_name"},
		}

		result := builder.Arguments(variables, nil)

		assert.Equal(t, "default_name", result["name"], "should use default value")
	})
}

func TestChatInputBuilder_PromptArguments(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	t.Run("basic prompt arguments", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "name", DefaultValue: "John"},
			{Name: "city", DefaultValue: "New York"},
		}

		result := builder.PromptArguments(variables)

		assert.Equal(t, "John", result["name"], "name should match default")
		assert.Equal(t, "New York", result["city"], "city should match default")
	})

	t.Run("empty variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}

		result := builder.PromptArguments(variables)

		assert.Empty(t, result, "result should be empty")
	})

	t.Run("nil variables", func(t *testing.T) {
		result := builder.PromptArguments(nil)

		assert.Empty(t, result, "result should be empty for nil input")
	})

	t.Run("variable with empty default", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "optional", DefaultValue: ""},
		}

		result := builder.PromptArguments(variables)

		assert.Equal(t, "", result["optional"], "empty default should be preserved")
	})
}

func TestChatInputBuilder_Options(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger)

	t.Run("add options to nil map", func(t *testing.T) {
		opts := map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  100,
		}

		result := builder.Options(opts, nil)

		assert.NotNil(t, result, "result should not be nil")
		assert.Len(t, result, 2, "should have two options")
	})

	t.Run("add options to existing map", func(t *testing.T) {
		existingVal, _ := structpb.NewValue("existing")
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"existing_key": existingAny,
		}

		opts := map[string]interface{}{
			"temperature": 0.7,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 2, "should have two options")
		assert.NotNil(t, result["existing_key"], "existing key should be preserved")
		assert.NotNil(t, result["temperature"], "new key should be added")
	})

	t.Run("override existing option", func(t *testing.T) {
		existingVal, _ := structpb.NewValue(0.5)
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"temperature": existingAny,
		}

		opts := map[string]interface{}{
			"temperature": 0.9,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should have one option")
		// The value should be overridden
	})

	t.Run("empty opts", func(t *testing.T) {
		existingVal, _ := structpb.NewValue("existing")
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"existing_key": existingAny,
		}

		opts := map[string]interface{}{}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should preserve existing options")
	})

	t.Run("various value types", func(t *testing.T) {
		opts := map[string]interface{}{
			"string_val": "hello",
			"int_val":    42,
			"float_val":  3.14,
			"bool_val":   true,
			"list_val":   []interface{}{1, 2, 3},
			"map_val":    map[string]interface{}{"nested": "value"},
		}

		result := builder.Options(opts, nil)

		assert.Len(t, result, 6, "should handle all value types")
	})
}

func TestChatInputBuilder_Integration(t *testing.T) {
	logger := newTestLogger()
	builder := NewChatInputBuilder(logger).(*inputChatBuilder)

	t.Run("full chat request workflow", func(t *testing.T) {
		// Create credential
		structVal, _ := structpb.NewStruct(map[string]interface{}{
			"api_key": "sk-test-key",
			"org_id":  "org-123",
		})
		credential := builder.Credential(42, structVal)

		// Create model options
		modelOpts := builder.Options(map[string]interface{}{
			"model.name":        "gpt-4",
			"model.temperature": 0.7,
			"model.max_tokens":  1000,
		}, nil)

		// Create tools
		tools := []*protos.FunctionDefinition{
			{
				Name:        "get_current_time",
				Description: "Get the current time in a specific timezone",
			},
		}

		// Create conversations
		conversations := []*protos.Message{
			builder.WithinMessage("system", "You are a helpful assistant."),
			builder.WithinMessage("user", "What time is it in Tokyo?"),
		}

		// Create additional data
		additionalData := map[string]string{
			"trace_id":   "trace-abc123",
			"session_id": "session-xyz789",
		}

		// Build the request
		request := builder.Chat(
			"request-001",
			credential,
			modelOpts,
			tools,
			additionalData,
			conversations...,
		)

		// Assertions
		assert.NotNil(t, request)
		assert.Equal(t, "request-001", request.RequestId)
		assert.Equal(t, uint64(42), request.Credential.Id)
		assert.Len(t, request.Conversations, 2)
		assert.Len(t, request.ToolDefinitions, 1)
		assert.Equal(t, "trace-abc123", request.AdditionalData["trace_id"])
		assert.NotNil(t, request.ModelParameters["model.name"])
	})
}
