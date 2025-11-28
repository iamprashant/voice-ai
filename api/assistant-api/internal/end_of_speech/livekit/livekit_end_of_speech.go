package internal_livekit_end_of_speech

import (
	internal_end_of_speech "github.com/rapidaai/api/assistant-api/internal/end_of_speech"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

func NewLivekitEndOfSpeech(
	logger commons.Logger,
	onCallback internal_end_of_speech.EndOfSpeechCallback,
	opts utils.Option,
) (internal_end_of_speech.EndOfSpeech, error) {
	return nil, nil
}
