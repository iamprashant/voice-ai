package internal_adapter_request_streamers

import (
	"context"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	lexatic_backend "github.com/rapidaai/protos"
	"google.golang.org/grpc"
)

type unidirectionalStreamer struct {
	server grpc.BidiStreamingServer[lexatic_backend.AssistantMessagingRequest, lexatic_backend.AssistantMessagingResponse]
}

func NewGrpcUnidirectionalStreamer(
	server lexatic_backend.TalkService_AssistantTalkServer) Streamer {
	return &unidirectionalStreamer{
		server: server,
	}
}

func (uds *unidirectionalStreamer) Context() context.Context {
	return uds.server.Context()
}

func (uds *unidirectionalStreamer) Recv() (*lexatic_backend.AssistantMessagingRequest, error) {
	return uds.server.Recv()
}

// Send sends an output value to the stream.
// It returns an error if the send operation fails.
func (uds *unidirectionalStreamer) Send(out *lexatic_backend.AssistantMessagingResponse) error {
	return uds.server.Send(out)
}

func (uds *unidirectionalStreamer) Config() *StreamAttribute {
	return &StreamAttribute{
		inputConfig: &StreamConfig{
			audio: internal_audio.NewLinear16khzMonoAudioConfig(),
			text: &struct {
				Charset string `json:"charset"`
			}{
				Charset: "UTF-8",
			},
		},
		outputConfig: &StreamConfig{
			audio: internal_audio.NewLinear16khzMonoAudioConfig(),
			text: &struct {
				Charset string `json:"charset"`
			}{
				Charset: "UTF-8",
			},
		},
	}
}
