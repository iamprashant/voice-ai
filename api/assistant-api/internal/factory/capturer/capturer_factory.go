package internal_capturer_factory

import (
	internal_capturers "github.com/rapidaai/api/assistant-api/internal/capturers"
	"github.com/rapidaai/pkg/commons"
)

type CapturerIdentifier string

const (
	Azure  CapturerIdentifier = "azure-cloud"
	Google CapturerIdentifier = "google-cloud"
	AWSS3  CapturerIdentifier = "aws-s3"
)

func GetAudioCapturer(aa CapturerIdentifier, logger commons.Logger, opts *internal_capturers.CapturerOptions) (internal_capturers.AudioCapturer, error) {
	switch aa {
	default:
		return internal_capturers.NewS3AudioCapturer(logger, opts)
	}
}

func GetTextCapturer(aa CapturerIdentifier, logger commons.Logger, opts *internal_capturers.CapturerOptions) (internal_capturers.TextCapturer, error) {
	switch aa {
	default:
		return internal_capturers.NewS3TextCapturer(logger, opts)
	}
}
