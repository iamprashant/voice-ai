// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter_generic

import (
	"context"
	"fmt"
	"sync"
	"time"

	internal_sentence_aggregator "github.com/rapidaai/api/assistant-api/internal/aggregator/text"
	internal_denoiser "github.com/rapidaai/api/assistant-api/internal/denoiser"
	internal_end_of_speech "github.com/rapidaai/api/assistant-api/internal/end_of_speech"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	internal_vad "github.com/rapidaai/api/assistant-api/internal/vad"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"golang.org/x/sync/errgroup"
)

func (io *GenericRequestor) Input(message *protos.ConversationUserMessage) error {
	switch msg := message.GetMessage().(type) {
	case *protos.ConversationUserMessage_Audio:
		return io.OnPacket(io.Context(), internal_type.UserAudioPacket{Audio: msg.Audio})
	case *protos.ConversationUserMessage_Text:
		return io.OnPacket(io.Context(), internal_type.UserTextPacket{Text: msg.Text})
	default:
		return fmt.Errorf("illegal input from the user %+v", msg)
	}

}

// Init initializes the audio talking system for a given assistant persona.
// It sets up both audio input and output transformer.
// This function is typically called at the beginning of a communication session.
func (listening *GenericRequestor) connectMicrophone(ctx context.Context, audioConfig *protos.AudioConfig) error {
	eGroup, ctx := errgroup.WithContext(ctx)
	options := utils.Option{"microphone.eos.timeout": 500}
	transformerConfig, _ := listening.GetSpeechToTextTransformer()
	if transformerConfig != nil && audioConfig != nil {
		options = utils.MergeMaps(options, transformerConfig.GetOptions())
		eGroup.Go(func() error {
			err := listening.initializeSpeechToText(ctx, transformerConfig, audioConfig, options)
			if err != nil {
				listening.logger.Errorf("unable to initialize transformer %+v", err)
			}
			return nil
		})

		eGroup.Go(func() error {
			err := listening.initializeVAD(ctx, audioConfig, options)
			if err != nil {
				listening.logger.Errorf("illegal input audio transformer, check the config and re-init")
			}
			return nil
		})

		eGroup.Go(func() error {
			err := listening.initializeDenoiser(ctx, audioConfig, options)
			if err != nil {
				listening.logger.Errorf("illegal input audio transformer, check the config and re-init")
			}
			return nil
		})

	}

	eGroup.Go(func() error {
		err := listening.initializeEndOfSpeech(ctx, options)
		if err != nil {
			listening.logger.Errorf("illegal input audio transformer, check the config and re-init")
		}
		return nil
	})

	if err := eGroup.Wait(); err != nil {
		return err
	}
	return nil
}

func (listening *GenericRequestor) disconnectMicrophone(ctx context.Context) error {
	if listening.speechToTextTransformer != nil {
		if err := listening.speechToTextTransformer.Close(ctx); err != nil {
			listening.logger.Warnf("cancel all output transformer with error %v", err)
		}
	}

	if listening.endOfSpeech != nil {
		if err := listening.endOfSpeech.Close(); err != nil {
			listening.logger.Warnf("cancel end of speech with error %v", err)
		}
	}

	if listening.vad != nil {
		if err := listening.vad.Close(); err != nil {
			listening.logger.Warnf("cancel vad with error %v", err)
		}
	}
	return nil
}

func (listening *GenericRequestor) initializeSpeechToText(ctx context.Context, transformerConfig *internal_assistant_entity.AssistantDeploymentAudio, audioConfig *protos.AudioConfig, options utils.Option) error {
	ctx, span, _ := listening.Tracer().StartSpan(ctx, utils.AssistantListenConnectStage)
	defer span.EndSpan(ctx, utils.AssistantListenConnectStage)

	span.AddAttributes(ctx,
		internal_telemetry.KV{K: "options", V: internal_telemetry.JSONValue(options)},
		internal_telemetry.KV{K: "provider", V: internal_telemetry.StringValue(transformerConfig.AudioProvider)},
	)

	credentialId, err := options.GetUint64("rapida.credential_id")
	if err != nil {
		listening.logger.Errorf("unable to find credential from options %+v", err)
		return err
	}
	credential, err := listening.VaultCaller().GetCredential(ctx, listening.Auth(), credentialId)
	if err != nil {
		listening.logger.Errorf("Api call to find credential failed %+v", err)
		return err
	}

	atransformer, err := internal_transformer.GetSpeechToTextTransformer(
		listening.Context(),
		listening.logger,
		transformerConfig.AudioProvider,
		credential,
		audioConfig,
		func(pkt ...internal_type.Packet) error { return listening.OnPacket(ctx, pkt...) },
		options)
	if err != nil {
		listening.logger.Errorf("unable to create input audio transformer with error %v", err)
		return err
	}
	err = atransformer.Initialize()
	if err != nil {
		listening.logger.Errorf("unable to initilize transformer %v", err)
		return err
	}
	listening.speechToTextTransformer = atransformer
	return nil
}

func (listening *GenericRequestor) initializeEndOfSpeech(ctx context.Context, options utils.Option) error {
	start := time.Now()

	endOfSpeech, err := internal_end_of_speech.GetEndOfSpeech(ctx,
		listening.logger,
		listening.OnPacket,
		options)
	if err != nil {
		listening.logger.Warnf("unable to initialize text analyzer %+v", err)
		return err
	}
	listening.endOfSpeech = endOfSpeech
	listening.logger.Benchmark("listen.endOfSpeech", time.Since(start))
	return nil
}

func (listening *GenericRequestor) initializeDenoiser(ctx context.Context, audioConfig *protos.AudioConfig, options utils.Option) error {
	denoise, err := internal_denoiser.GetDenoiser(listening.Context(), listening.logger, audioConfig, options)
	if err != nil {
		listening.logger.Errorf("error wile intializing denoiser %+v", err)
	}
	listening.denoiser = denoise
	return nil
}

func (listening *GenericRequestor) initializeVAD(ctx context.Context, audioConfig *protos.AudioConfig, options utils.Option,
) error {
	start := time.Now()
	vad, err := internal_vad.GetVAD(listening.Context(), listening.logger, audioConfig, listening.OnPacket, options)
	if err != nil {
		listening.logger.Errorf("error wile intializing vad %+v", err)
		return err
	}
	listening.vad = vad
	listening.logger.Benchmark("listen.initializeVAD", time.Since(start))
	return nil
}

// Init initializes the audio talking system for a given assistant persona.
// It sets up both audio input and output transformer.
// This function is typically called at the beginning of a communication session.
func (spk *GenericRequestor) connectSpeaker(ctx context.Context, audioOutConfig *protos.AudioConfig) error {
	speakerOpts := spk.GetOptions()
	var wg sync.WaitGroup
	// initialize audio output transformer
	outputTransformer, _ := spk.GetTextToSpeechTransformer()
	if outputTransformer != nil && audioOutConfig != nil {
		speakerOpts = utils.MergeMaps(outputTransformer.GetOptions())
		wg.Add(1)
		utils.Go(ctx, func() {
			defer wg.Done()
			if err := spk.initializeTextToSpeech(ctx, outputTransformer, audioOutConfig, speakerOpts); err != nil {
				spk.logger.Errorf("unable to initialize text to speech transformer with error %v", err)
				return
			}
		})
	}
	//

	wg.Add(1)
	utils.Go(ctx, func() {
		defer wg.Done()
		if err := spk.initializeTextAggregator(ctx, speakerOpts); err != nil {
			spk.logger.Errorf("unable to initialize sentence assembler with error %v", err)
			return
		}
	})

	wg.Wait()
	return nil
}

func (spk *GenericRequestor) disconnectSpeaker() error {
	if spk.textToSpeechTransformer != nil {
		if err := spk.textToSpeechTransformer.Close(spk.Context()); err != nil {
			spk.logger.Errorf("cancel all output transformer with error %v", err)
		}
	}
	if spk.textAggregator != nil {
		spk.textAggregator.Close()
	}
	return nil
}

func (spk *GenericRequestor) initializeTextToSpeech(context context.Context, transformerConfig *internal_assistant_entity.AssistantDeploymentAudio, audioConfig *protos.AudioConfig, speakerOpts utils.Option) error {
	context, span, _ := spk.Tracer().StartSpan(context, utils.AssistantSpeakConnectStage)
	defer span.EndSpan(context, utils.AssistantSpeakConnectStage)
	span.AddAttributes(context,
		internal_adapter_telemetry.KV{
			K: "options", V: internal_adapter_telemetry.JSONValue(speakerOpts),
		},
		internal_adapter_telemetry.KV{
			K: "provider", V: internal_adapter_telemetry.StringValue(transformerConfig.AudioProvider),
		},
	)
	credentialId, err := speakerOpts.GetUint64("rapida.credential_id")
	if err != nil {
		spk.logger.Errorf("unable to find credential from options %+v", err)
		return err
	}
	credential, err := spk.VaultCaller().GetCredential(context, spk.Auth(), credentialId)
	if err != nil {
		spk.logger.Errorf("Api call to find credential failed %+v", err)
		return err
	}

	atransformer, err := internal_transformer.GetTextToSpeechTransformer(
		context, spk.logger,
		transformerConfig.GetName(),
		credential, audioConfig,
		func(pkt ...internal_type.Packet) error { return spk.OnPacket(context, pkt...) },
		speakerOpts)
	if err != nil {
		spk.logger.Errorf("unable to create input audio transformer with error %v", err)
		return err
	}
	if err := atransformer.Initialize(); err != nil {
		spk.logger.Errorf("unable to initilize transformer %v", err)
		return err
	}
	spk.textToSpeechTransformer = atransformer
	return nil
}

func (spk *GenericRequestor) initializeTextAggregator(ctx context.Context, options utils.Option) error {
	if textAggregator, err := internal_sentence_aggregator.GetLLMTextAggregator(spk.Context(), spk.logger, options); err == nil {
		spk.textAggregator = textAggregator
		go spk.onAssembleSentence(spk.Context())
	}
	return nil
}

func (spk *GenericRequestor) onAssembleSentence(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			spk.logger.Debugf("OnCompleteSentence stopped due to context cancellation")
			return

		case result, ok := <-spk.textAggregator.Result():
			if !ok {
				spk.logger.Debugf("speak: OnCompleteSentence tokenizer channel closed")
				return
			}
			spk.callSpeaking(ctx, result)
		}
	}
}
