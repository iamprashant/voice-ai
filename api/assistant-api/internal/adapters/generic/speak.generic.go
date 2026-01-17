// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter_generic

import (
	"context"
	"sync"
	"time"

	internal_sentence_assembler "github.com/rapidaai/api/assistant-api/internal/assembler/sentence"
	internal_synthesizers "github.com/rapidaai/api/assistant-api/internal/synthesizes"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Init initializes the audio talking system for a given assistant persona.
// It sets up both audio input and output transformer.
// This function is typically called at the beginning of a communication session.
func (spk *GenericRequestor) ConnectSpeaker(ctx context.Context, audioOutConfig *protos.AudioConfig) error {
	context, span, _ := spk.Tracer().StartSpan(ctx, utils.AssistantSpeakConnectStage)
	defer span.EndSpan(context, utils.AssistantSpeakConnectStage)

	// speaker options
	speakerOpts := spk.GetOptions()
	start := time.Now()
	var wg sync.WaitGroup

	// initialize audio output transformer
	if audioOutConfig != nil {
		if outputTransformer, err := spk.GetTextToSpeechTransformer(); err == nil {
			speakerOpts = utils.MergeMaps(outputTransformer.GetOptions())

			span.AddAttributes(context,
				internal_adapter_telemetry.KV{
					K: "options", V: internal_adapter_telemetry.JSONValue(speakerOpts),
				},
				internal_adapter_telemetry.KV{
					K: "provider", V: internal_adapter_telemetry.StringValue(outputTransformer.AudioProvider),
				},
			)
			//
			wg.Add(1)
			utils.Go(context, func() {
				defer wg.Done()
				opts := &internal_type.TextToSpeechInitializeOptions{
					AudioConfig:  audioOutConfig,
					OnSpeech:     func(pkt ...internal_type.Packet) error { return spk.OnPacket(context, pkt...) },
					ModelOptions: speakerOpts,
				}

				credentialId, err := opts.ModelOptions.GetUint64("rapida.credential_id")
				if err != nil {
					spk.logger.Errorf("unable to find credential from options %+v", err)
					return
				}
				credential, err := spk.VaultCaller().GetCredential(context, spk.Auth(), credentialId)
				if err != nil {
					spk.logger.Errorf("Api call to find credential failed %+v", err)
					return
				}

				atransformer, err := internal_transformer.GetTextToSpeechTransformer(internal_transformer.AudioTransformer(outputTransformer.GetName()), context, spk.logger, credential, opts)
				if err != nil {
					spk.logger.Errorf("unable to create input audio transformer with error %v", err)
					return
				}
				spk.logger.Benchmark("speak.transformer.GetOutputAudioTransformer", time.Since(start))
				if err := atransformer.Initialize(); err != nil {
					spk.logger.Errorf("unable to initilize transformer %v", err)
					return
				}
				spk.textToSpeechTransformer = atransformer
				spk.logger.Benchmark("speak.transformer.Initialize", time.Since(start))
			})
		}
	}
	//

	wg.Add(1)
	utils.Go(context, func() {
		defer wg.Done()
		if sentenceAssembler, err := internal_sentence_assembler.NewLLMSentenceAssembler(spk.logger, speakerOpts); err == nil {
			spk.sentenceAssembler = sentenceAssembler
			go spk.OnCompleteSentence(context)
		}
		if normalizer, err := internal_synthesizers.NewSentenceNormalizeSynthesizer(spk.logger, internal_synthesizers.SynthesizerOptions{SpeakerOptions: speakerOpts}); err == nil {
			spk.synthesizers = append(spk.synthesizers, normalizer)
		}
		if formatter, err := internal_synthesizers.NewSentenceFormattingSynthesizer(spk.logger, internal_synthesizers.SynthesizerOptions{SpeakerOptions: speakerOpts}); err == nil {
			spk.synthesizers = append(spk.synthesizers, formatter)
		}
		spk.logger.Benchmark("speak.GetAudioOutputTransformer.synthesizers", time.Since(start))
	})

	wg.Wait()
	spk.logger.Benchmark("speak.Init", time.Since(start))
	return nil
}

func (spk *GenericRequestor) OnCompleteSentence(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			spk.logger.Debugf("OnCompleteSentence stopped due to context cancellation")
			return

		case result, ok := <-spk.sentenceAssembler.Result():
			if !ok {
				spk.logger.Debugf("speak: OnCompleteSentence tokenizer channel closed")
				return
			}
			//
			switch res := result.(type) {
			case internal_type.LLMMessagePacket:
				ctx, span, _ := spk.Tracer().StartSpan(spk.Context(), utils.AssistantSpeakingStage)
				defer span.EndSpan(ctx, utils.AssistantSpeakingStage)
				span.AddAttributes(ctx,
					internal_adapter_telemetry.MessageKV(res.ContextID),
					internal_adapter_telemetry.KV{K: "activity", V: internal_adapter_telemetry.StringValue("finish_speaking")},
				)
				if spk.textToSpeechTransformer != nil {
					if err := spk.textToSpeechTransformer.Transform(spk.Context(), res); err != nil {
						spk.logger.Errorf("speak: failed to send flush to text to speech transformer error: %v", err)
					}
				}
			case internal_type.LLMStreamPacket:
				ctxSpan, span, _ := spk.Tracer().StartSpan(ctx, utils.AssistantSpeakingStage)
				span.AddAttributes(ctxSpan,
					internal_adapter_telemetry.MessageKV(res.ContextID),
					internal_adapter_telemetry.KV{
						K: "activity", V: internal_adapter_telemetry.StringValue("speak"),
					},
					internal_adapter_telemetry.KV{
						K: "script", V: internal_adapter_telemetry.StringValue(res.Text),
					},
				)

				for _, v := range spk.synthesizers {
					res = v.Synthesize(spk.Context(), res)
				}

				span.AddAttributes(ctxSpan,
					internal_adapter_telemetry.KV{
						K: "synthesize_script",
						V: internal_adapter_telemetry.StringValue(res.Text),
					},
				)

				// if err := spk.messaging.Transition(internal_adapter_request_customizers.AgentSpeaking); err == nil {
				if spk.textToSpeechTransformer != nil {
					if err := spk.textToSpeechTransformer.Transform(spk.Context(), res); err != nil {
						spk.logger.Errorf("speak: failed to send sentence to text to speech transformer error: %v", err)
					}
				}
				if err := spk.Notify(ctx, &protos.AssistantConversationAssistantMessage{Time: timestamppb.Now(), Id: res.ContextId(), Completed: true, Message: &protos.AssistantConversationAssistantMessage_Text{Text: &protos.AssistantConversationMessageTextContent{Content: res.Text}}}); err != nil {
					spk.logger.Tracef(ctx, "error while outputting chunk to the user: %w", err)
				}
				// }

				span.EndSpan(ctxSpan, utils.AssistantSpeakingStage)
			default:
			}

		}
	}
}

func (spk *GenericRequestor) CloseSpeaker() error {
	if spk.textToSpeechTransformer != nil {
		if err := spk.textToSpeechTransformer.Close(spk.Context()); err != nil {
			spk.logger.Errorf("cancel all output transformer with error %v", err)
		}
	}
	if spk.sentenceAssembler != nil {
		spk.sentenceAssembler.Close()
	}
	return nil
}
