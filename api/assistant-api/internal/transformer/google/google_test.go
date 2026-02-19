package internal_transformer_google

import (
	"testing"

	"cloud.google.com/go/speech/apiv2/speechpb"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
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

func TestNewGoogleOption_ValidCredentials(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":                 "test-api-key",
		"project_id":          "test-project",
		"service_account_key": `{"type":"service_account"}`,
	})
	opt, err := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Equal(t, "test-project", opt.projectId)
	assert.Len(t, opt.clientOptons, 3) // API key + quota project + credentials JSON
}

func TestNewGoogleOption_OnlyAPIKey(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key": "test-api-key",
	})
	opt, err := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Len(t, opt.clientOptons, 1)
}

func TestNewGoogleOption_EmptyVault(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{})
	opt, err := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err) // Google constructor doesn't error on missing keys
	assert.NotNil(t, opt)
	assert.Empty(t, opt.clientOptons)
}

func TestNewGoogleOption_EmptyStringValues(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{
		"key":        "",
		"project_id": "",
	})
	opt, err := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	assert.NoError(t, err)
	assert.NotNil(t, opt)
	assert.Empty(t, opt.clientOptons) // Empty strings are skipped
}

// --- SpeechToTextOptions Tests ---

func TestSpeechToTextOptions_Defaults(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opt, _ := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	sttOpts := opt.SpeechToTextOptions()

	assert.NotNil(t, sttOpts)
	assert.NotNil(t, sttOpts.Config)

	// Verify encoding config
	decodingCfg := sttOpts.Config.GetExplicitDecodingConfig()
	assert.NotNil(t, decodingCfg)
	assert.Equal(t, speechpb.ExplicitDecodingConfig_LINEAR16, decodingCfg.Encoding)
	assert.Equal(t, int32(16000), decodingCfg.SampleRateHertz)
	assert.Equal(t, int32(1), decodingCfg.AudioChannelCount)

	// Verify features
	assert.True(t, sttOpts.Config.Features.EnableAutomaticPunctuation)
	assert.True(t, sttOpts.Config.Features.EnableWordConfidence)
	assert.True(t, sttOpts.Config.Features.ProfanityFilter)
	assert.True(t, sttOpts.Config.Features.EnableSpokenPunctuation)

	// Verify default language
	assert.Equal(t, []string{DefaultLanguageCode}, sttOpts.Config.LanguageCodes)

	// Verify model
	assert.Equal(t, "long", sttOpts.Config.Model)

	// Verify streaming features
	assert.False(t, sttOpts.StreamingFeatures.EnableVoiceActivityEvents)
	assert.True(t, sttOpts.StreamingFeatures.InterimResults)
}

func TestSpeechToTextOptions_WithLanguageOverride(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"listen.language": "fr-FR",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	sttOpts := opt.SpeechToTextOptions()

	assert.Equal(t, []string{"fr-FR"}, sttOpts.Config.LanguageCodes)
}

func TestSpeechToTextOptions_WithMultipleLanguages(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"listen.language": "en-US" + commons.SEPARATOR + "fr-FR" + commons.SEPARATOR + "de-DE",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	sttOpts := opt.SpeechToTextOptions()

	assert.Equal(t, []string{"en-US", "fr-FR", "de-DE"}, sttOpts.Config.LanguageCodes)
}

func TestSpeechToTextOptions_WithModelOverride(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"listen.model": "chirp",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	sttOpts := opt.SpeechToTextOptions()

	assert.Equal(t, "chirp", sttOpts.Config.Model)
}

// --- TextToSpeechOptions Tests ---

func TestTextToSpeechOptions_Defaults(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opt, _ := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	ttsOpts := opt.TextToSpeechOptions()

	assert.NotNil(t, ttsOpts)
	assert.Equal(t, DefaultVoice, ttsOpts.Voice.Name)
	assert.Equal(t, texttospeechpb.AudioEncoding_PCM, ttsOpts.StreamingAudioConfig.AudioEncoding)
	assert.Equal(t, int32(16000), ttsOpts.StreamingAudioConfig.SampleRateHertz)
}

func TestTextToSpeechOptions_WithVoiceOverride(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"speak.voice.id": "en-US-Wavenet-D",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	ttsOpts := opt.TextToSpeechOptions()

	assert.Equal(t, "en-US-Wavenet-D", ttsOpts.Voice.Name)
	// Encoding still hardcoded
	assert.Equal(t, texttospeechpb.AudioEncoding_PCM, ttsOpts.StreamingAudioConfig.AudioEncoding)
	assert.Equal(t, int32(16000), ttsOpts.StreamingAudioConfig.SampleRateHertz)
}

// --- GetRecognizer Tests ---

func TestGetRecognizer_Default(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "my-project"})
	opt, _ := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	recognizer := opt.GetRecognizer()

	assert.Equal(t, "projects/my-project/locations/global/recognizers/_", recognizer)
}

func TestGetRecognizer_WithGlobalRegion(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "my-project"})
	opts := utils.Option{
		"listen.region": "global",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	recognizer := opt.GetRecognizer()

	assert.Equal(t, "projects/my-project/locations/global/recognizers/_", recognizer)
}

func TestGetRecognizer_WithSpecificRegion(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "my-project"})
	opts := utils.Option{
		"listen.region": "us-central1",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	recognizer := opt.GetRecognizer()

	assert.Equal(t, "projects/my-project/locations/us-central1/recognizers/_", recognizer)
}

// --- GetSpeechToTextClientOptions Tests ---

func TestGetSpeechToTextClientOptions_Default(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opt, _ := NewGoogleOption(newTestLogger(), cred, utils.Option{})
	clientOpts := opt.GetSpeechToTextClientOptions()

	// Should return the base client options without additional endpoint
	assert.Equal(t, len(opt.clientOptons), len(clientOpts))
}

func TestGetSpeechToTextClientOptions_WithRegion(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"listen.region": "eu-west1",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	clientOpts := opt.GetSpeechToTextClientOptions()

	// Should have one more option (endpoint) than base client options
	assert.Equal(t, len(opt.clientOptons)+1, len(clientOpts))
}

func TestGetSpeechToTextClientOptions_GlobalRegion(t *testing.T) {
	cred := newVaultCredential(map[string]interface{}{"key": "k", "project_id": "p"})
	opts := utils.Option{
		"listen.region": "global",
	}
	opt, _ := NewGoogleOption(newTestLogger(), cred, opts)
	clientOpts := opt.GetSpeechToTextClientOptions()

	// Global region should NOT add endpoint override
	assert.Equal(t, len(opt.clientOptons), len(clientOpts))
}
