package internal_transformer_factory

import (
	"context"
	"fmt"

	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	internal_transformer_assemblyai "github.com/rapidaai/api/assistant-api/internal/transformer/assembly-ai"
	internal_transformer_aws "github.com/rapidaai/api/assistant-api/internal/transformer/aws"
	internal_transformer_azure "github.com/rapidaai/api/assistant-api/internal/transformer/azure"
	internal_transformer_cartesia "github.com/rapidaai/api/assistant-api/internal/transformer/cartesia"
	internal_transformer_deepgram "github.com/rapidaai/api/assistant-api/internal/transformer/deepgram"
	internal_transformer_elevenlabs "github.com/rapidaai/api/assistant-api/internal/transformer/elevenlabs"
	internal_transformer_google "github.com/rapidaai/api/assistant-api/internal/transformer/google"
	internal_transformer_revai "github.com/rapidaai/api/assistant-api/internal/transformer/revai"
	internal_transformer_sarvam "github.com/rapidaai/api/assistant-api/internal/transformer/sarvam"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type AudioTransformer string

const (
	DEEPGRAM     AudioTransformer = "deepgram"
	AZURE        AudioTransformer = "azure"
	CARTESIA     AudioTransformer = "cartesia"
	GOOGLE       AudioTransformer = "google"
	GOOGLE_CLOUD AudioTransformer = "google-cloud"
	AWS          AudioTransformer = "aws"
	REVAI        AudioTransformer = "revai"
	SARVAM       AudioTransformer = "sarvamai"
	ELEVENLABS   AudioTransformer = "elevenlabs"
	ASSEMBLYAI   AudioTransformer = "assemblyai"
)

func (at AudioTransformer) String() string {
	return string(at)
}
func GetTextToSpeechTransformer(
	at AudioTransformer,
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	opts *internal_transformer.TextToSpeechInitializeOptions) (internal_transformer.TextToSpeechTransformer, error) {

	switch at {
	case DEEPGRAM:
		return internal_transformer_deepgram.NewDeepgramTextToSpeech(ctx, logger, credential, opts)
	case AZURE:
		return internal_transformer_azure.NewAzureTextToSpeech(ctx, logger, credential, opts)
	case CARTESIA:
		return internal_transformer_cartesia.NewCartesiaTextToSpeech(ctx, logger, credential, opts)
	case GOOGLE, GOOGLE_CLOUD:
		return internal_transformer_google.NewGoogleTextToSpeech(ctx, logger, credential, opts)
	case AWS:
		return internal_transformer_aws.NewAWSTextToSpeech(ctx, logger, credential, opts)
	case REVAI:
		return internal_transformer_revai.NewRevaiTextToSpeech(ctx, logger, credential, opts)
	case SARVAM:
		return internal_transformer_sarvam.NewSarvamTextToSpeech(ctx, logger, credential, opts)
	case ELEVENLABS:
		return internal_transformer_elevenlabs.NewElevenlabsTextToSpeech(ctx, logger, credential, opts)
	default:
		return nil, fmt.Errorf("illegal text to speech idenitfier")
	}
}

func GetSpeechToTextTransformer(at AudioTransformer,
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	opts *internal_transformer.SpeechToTextInitializeOptions) (internal_transformer.SpeechToTextTransformer, error) {

	switch at {
	case DEEPGRAM:
		return internal_transformer_deepgram.NewDeepgramSpeechToText(ctx, logger, credential, opts)
	case AZURE:
		return internal_transformer_azure.NewAzureSpeechToText(ctx, logger, credential, opts)
	case GOOGLE, GOOGLE_CLOUD:
		return internal_transformer_google.NewGoogleSpeechToText(ctx, logger, credential, opts)
	case AWS:
		return internal_transformer_aws.NewAWSSpeechToText(ctx, logger, credential, opts)
	case ASSEMBLYAI:
		return internal_transformer_assemblyai.NewAssemblyaiSpeechToText(ctx, logger, credential, opts)
	case REVAI:
		return internal_transformer_revai.NewRevaiSpeechToText(ctx, logger, credential, opts)
	case SARVAM:
		return internal_transformer_sarvam.NewSarvamSpeechToText(ctx, logger, credential, opts)
	case CARTESIA:
		return internal_transformer_cartesia.NewCartesiaSpeechToText(ctx, logger, credential, opts)
	default:
		return nil, fmt.Errorf("illegal speech to text idenitfier")
	}
}
