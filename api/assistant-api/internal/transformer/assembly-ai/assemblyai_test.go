package internal_transformer_assemblyai

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

func TestNewAssemblyaiOption_ValidCredentials(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "test-api-key"})
	opt, err := NewAssemblyaiOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Equal(t, "test-api-key", opt.GetKey())
}

func TestNewAssemblyaiOption_MissingKey(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"other": "value"})
	opt, err := NewAssemblyaiOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
	assert.Contains(t, err.Error(), "illegal vault config")
}

func TestNewAssemblyaiOption_EmptyVault(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{})
	opt, err := NewAssemblyaiOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
}

// --- Encoding Tests ---

func TestAssemblyaiGetEncoding(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k"})
	opt, _ := NewAssemblyaiOption(newTestLogger(), cred, utils.Option{})
	assert.Equal(t, "pcm_s16le", opt.GetEncoding())
}

// --- GetSpeechToTextConnectionString Tests ---

func TestGetSpeechToTextConnectionString_Default(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k"})
	opt, _ := NewAssemblyaiOption(newTestLogger(), cred, utils.Option{})
	connStr := opt.GetSpeechToTextConnectionString()

	assert.Contains(t, connStr, "wss://streaming.assemblyai.com/v3/ws?")
	assert.Contains(t, connStr, "sample_rate=16000")
	assert.Contains(t, connStr, "encoding=pcm_s16le")
	assert.Contains(t, connStr, "format_turns=true")
	assert.NotContains(t, connStr, "language=")
	assert.NotContains(t, connStr, "model=")
}

func TestGetSpeechToTextConnectionString_WithLanguage(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k"})
	opts := utils.Option{
		"listen.language": "fr",
	}
	opt, _ := NewAssemblyaiOption(newTestLogger(), cred, opts)
	connStr := opt.GetSpeechToTextConnectionString()

	assert.Contains(t, connStr, "language=fr")
	assert.Contains(t, connStr, "sample_rate=16000")
	assert.Contains(t, connStr, "encoding=pcm_s16le")
}

func TestGetSpeechToTextConnectionString_WithModel(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k"})
	opts := utils.Option{
		"listen.model": "nano",
	}
	opt, _ := NewAssemblyaiOption(newTestLogger(), cred, opts)
	connStr := opt.GetSpeechToTextConnectionString()

	assert.Contains(t, connStr, "model=nano")
	assert.Contains(t, connStr, "sample_rate=16000")
	assert.Contains(t, connStr, "encoding=pcm_s16le")
}

func TestGetSpeechToTextConnectionString_AllOptions(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k"})
	opts := utils.Option{
		"listen.language": "es",
		"listen.model":    "best",
	}
	opt, _ := NewAssemblyaiOption(newTestLogger(), cred, opts)
	connStr := opt.GetSpeechToTextConnectionString()

	assert.Contains(t, connStr, "language=es")
	assert.Contains(t, connStr, "model=best")
	assert.Contains(t, connStr, "sample_rate=16000")
	assert.Contains(t, connStr, "encoding=pcm_s16le")
	assert.Contains(t, connStr, "format_turns=true")
}
