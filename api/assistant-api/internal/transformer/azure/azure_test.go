package internal_transformer_azure

import (
	"testing"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
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

func TestNewAzureOption_ValidCredentials(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"subscription_key": "test-sub-key",
		"endpoint":         "https://test.cognitiveservices.azure.com",
	})
	opt, err := NewAzureOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Equal(t, "test-sub-key", opt.subscriptionKey)
	assert.Equal(t, "https://test.cognitiveservices.azure.com", opt.endpoint)
}

func TestNewAzureOption_MissingSubscriptionKey(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"endpoint": "https://test.cognitiveservices.azure.com",
	})
	opt, err := NewAzureOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
	assert.Contains(t, err.Error(), "subscription_key")
}

func TestNewAzureOption_MissingEndpoint(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"subscription_key": "test-sub-key",
	})
	opt, err := NewAzureOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
	assert.Contains(t, err.Error(), "endpoint")
}

func TestNewAzureOption_EmptyVault(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{})
	opt, err := NewAzureOption(newTestLogger(), cred, utils.Option{})
	assert.Error(t, err)
	assert.Nil(t, opt)
}

// --- Output Format Tests ---

func TestGetSpeechSynthesisOutputFormat(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"subscription_key": "k",
		"endpoint":         "https://e.azure.com",
	})
	opt, _ := NewAzureOption(newTestLogger(), cred, utils.Option{})
	assert.Equal(t, common.Raw16Khz16BitMonoPcm, opt.GetSpeechSynthesisOutputFormat())
}

// --- Audio Stream Format Tests ---

func TestGetAudioStreamFormat(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"subscription_key": "k",
		"endpoint":         "https://e.azure.com",
	})
	opt, _ := NewAzureOption(newTestLogger(), cred, utils.Option{})
	format := opt.GetAudioStreamFormat()
	assert.NotNil(t, format)
}
