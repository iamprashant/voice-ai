package internal_transformer_resemble

import (
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func newTestLogger() commons.Logger {
	l, _ := commons.NewApplicationLogger()
	return l
}

func newVaultCredential(m map[string]interface{}) *protos.VaultCredential {
	val, _ := structpb.NewStruct(m)
	return &protos.VaultCredential{Value: val}
}

// --- Constructor Tests ---

func TestNewResembleOption_ValidCredentials(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "test-api-key",
		"project_id": "test-project",
	})
	opt, err := NewResembleOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Equal(t, "test-api-key", opt.GetKey())
	assert.Equal(t, "test-project", opt.GetProject())
}

func TestNewResembleOption_MissingKey(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"project_id": "test-project",
	})
	opt, err := NewResembleOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
	assert.Contains(t, err.Error(), "illegal vault config")
}

func TestNewResembleOption_MissingProjectId(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key": "test-api-key",
	})
	opt, err := NewResembleOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
	assert.Contains(t, err.Error(), "illegal vault config")
}

func TestNewResembleOption_EmptyVault(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{})
	opt, err := NewResembleOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
}

// --- Encoding Tests ---

func TestResembleGetEncoding(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "k",
		"project_id": "p",
	})
	opt, _ := NewResembleOption(newTestLogger(), cred, utils.Option{})
	assert.Equal(t, "PCM_16", opt.GetEncoding())
}

// --- GetTextToSpeechRequest Tests ---

func TestGetTextToSpeechRequest_Default(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "k",
		"project_id": "my-project",
	})
	opt, _ := NewResembleOption(newTestLogger(), cred, utils.Option{})
	req := opt.GetTextToSpeechRequest("ctx-123", "Hello world")
	assert.Equal(t, VOICE_ID, req["voice_uuid"])
	assert.Equal(t, "ctx-123", req["request_id"])
	assert.Equal(t, "my-project", req["project_uuid"])
	assert.Equal(t, "Hello world", req["data"])
	assert.Equal(t, true, req["binary_response"])
	assert.Equal(t, "PCM_16", req["precision"])
	assert.Equal(t, 16000, req["sample_rate"])
}

func TestGetTextToSpeechRequest_DifferentContext(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "k",
		"project_id": "proj-456",
	})
	opt, _ := NewResembleOption(newTestLogger(), cred, utils.Option{})
	req := opt.GetTextToSpeechRequest("req-789", "Bonjour le monde")
	assert.Equal(t, "req-789", req["request_id"])
	assert.Equal(t, "proj-456", req["project_uuid"])
	assert.Equal(t, "Bonjour le monde", req["data"])
	assert.Equal(t, "PCM_16", req["precision"])
	assert.Equal(t, 16000, req["sample_rate"])
}

func TestGetTextToSpeechRequest_EmptyText(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "k",
		"project_id": "p",
	})
	opt, _ := NewResembleOption(newTestLogger(), cred, utils.Option{})
	req := opt.GetTextToSpeechRequest("ctx", "")
	assert.Equal(t, "", req["data"])
	assert.Equal(t, "PCM_16", req["precision"])
	assert.Equal(t, 16000, req["sample_rate"])
}
