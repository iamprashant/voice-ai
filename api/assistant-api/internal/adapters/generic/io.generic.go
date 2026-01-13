// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter_request_generic

import (
	"context"
	"fmt"
	"time"

	internal_adapter_request_customizers "github.com/rapidaai/api/assistant-api/internal/adapters/customizers"
	internal_end_of_speech "github.com/rapidaai/api/assistant-api/internal/end_of_speech"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Start Input
// =====================================
func (io *GenericRequestor) Input(message *protos.AssistantConversationUserMessage) error {
	switch msg := message.GetMessage().(type) {
	case *protos.AssistantConversationUserMessage_Audio:
		return io.InputAudio(io.Context(), msg.Audio.GetContent())
	case *protos.AssistantConversationUserMessage_Text:
		return io.InputText(io.Context(), msg.Text.GetContent())
	default:
		return fmt.Errorf("illegal input from the user %+v", msg)
	}

}

func (io *GenericRequestor) InputAudio(ctx context.Context, in []byte) error {
	if v, err := io.ListenAudio(ctx, in); err == nil {
		utils.Go(context.Background(), func() {
			io.recorder.User(v)
		})
	}
	return nil
}

func (io *GenericRequestor) InputText(ctx context.Context, msg string) error {
	// mark it interrupted
	io.messaging.Transition(internal_adapter_request_customizers.Interrupted)
	//
	interim := io.messaging.Create(type_enums.UserActor, msg)
	if err := io.Notify(ctx, &protos.AssistantConversationUserMessage{
		Id:        interim.GetId(),
		Completed: false,
		Message: &protos.AssistantConversationUserMessage_Text{
			Text: &protos.AssistantConversationMessageTextContent{
				Content: interim.String(),
			},
		},
		Time: timestamppb.Now(),
	}); err != nil {
		io.logger.Tracef(ctx, "error while notifying the text input from user: %w", err)
	}
	return io.ListenText(ctx, &internal_end_of_speech.UserEndOfSpeechInput{
		Message: interim.String(),
		Time:    time.Now(),
	},
	)
}
