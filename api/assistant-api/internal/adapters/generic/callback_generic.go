// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter_generic

import (
	"context"
	"errors"
	"time"

	internal_adapter_request_customizers "github.com/rapidaai/api/assistant-api/internal/adapters/customizers"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (talking *GenericRequestor) callEndOfSpeech(ctx context.Context, vl internal_type.Packet) error {
	if talking.endOfSpeech != nil {
		utils.Go(ctx, func() {
			if err := talking.endOfSpeech.Analyze(ctx, vl); err != nil {
				talking.logger.Errorf("end of speech analyze error: %v", err)
			}
		})
		return nil
	}
	return errors.New("end of speech analyzer not configured")
}

func (talking *GenericRequestor) callTextAggregator(ctx context.Context, vl internal_type.Packet) error {
	if talking.textAggregator != nil {
		if err := talking.textAggregator.Aggregate(ctx, vl); err != nil {
			talking.logger.Debugf("unable to send packet to aggregator %v", err)
		}
		return nil
	}
	return errors.New("textAggregator not configured")
}

func (talking *GenericRequestor) callRecording(ctx context.Context, vl internal_type.Packet) error {
	if talking.recorder != nil {
		utils.Go(ctx, func() {
			if err := talking.recorder.Record(ctx, vl); err != nil {
				talking.logger.Errorf("recorder error: %v", err)
			}
		})
		return nil
	}
	return nil
}

func (talking *GenericRequestor) callCreateMessage(ctx context.Context, vl internal_type.MessagePacket) error {
	utils.Go(ctx, func() {
		if err := talking.onCreateMessage(talking.Context(), vl); err != nil {
			talking.logger.Errorf("Error in onCreateMessage: %v", err)
		}
	})
	return nil
}

func (talking *GenericRequestor) callVadProcess(ctx context.Context, vl internal_type.UserAudioPacket) error {
	if talking.vad != nil {
		utils.Go(ctx, func() {
			if err := talking.vad.Process(ctx, vl); err != nil {
				talking.logger.Warnf("error while processing with vad %s", err.Error())
			}
		})
	}
	return nil
}

func (talking *GenericRequestor) callSpeechToText(ctx context.Context, vl internal_type.UserAudioPacket) error {
	if talking.speechToTextTransformer != nil {
		utils.Go(ctx, func() {
			if err := talking.speechToTextTransformer.Transform(ctx, vl); err != nil {
				talking.logger.Tracef(ctx, "error while transforming input %s and error %s", talking.speechToTextTransformer.Name(), err.Error())
			}
		})
	}
	return nil
}

func (spk *GenericRequestor) interruptAllProvider(ctx context.Context, result internal_type.InterruptionPacket) error {
	if spk.textToSpeechTransformer != nil {
		// can be done on goroutine
		utils.Go(ctx, func() {
			if err := spk.textToSpeechTransformer.Transform(spk.Context(), result); err != nil {
				spk.logger.Errorf("speak: failed to send interruption: %v", err)
			}
		})
	}

	if spk.assistantExecutor != nil {
		// can be done on goroutine
		utils.Go(ctx, func() {
			if err := spk.assistantExecutor.Execute(spk.Context(), spk, result); err != nil {
				spk.logger.Errorf("assistant executor interrupt error: %v", err)
			}
		})
	}
	return nil
}

func (spk *GenericRequestor) callSpeaking(ctx context.Context, result internal_type.LLMPacket) error {
	switch res := result.(type) {
	case internal_type.LLMResponseDonePacket:
		if spk.textToSpeechTransformer != nil {

			if result.ContextId() != spk.messaging.GetID() {
				return nil
			}
			ctx, span, _ := spk.Tracer().StartSpan(spk.Context(), utils.AssistantSpeakingStage)
			defer span.EndSpan(ctx, utils.AssistantSpeakingStage)
			span.AddAttributes(ctx,
				internal_adapter_telemetry.MessageKV(res.ContextID),
				internal_adapter_telemetry.KV{K: "activity", V: internal_adapter_telemetry.StringValue("finish_speaking")},
			)
			if err := spk.textToSpeechTransformer.Transform(spk.Context(), res); err != nil {
				spk.logger.Errorf("speak: failed to send flush to text to speech transformer error: %v", err)
			}
		}
	case internal_type.LLMResponseDeltaPacket:
		if result.ContextId() != spk.messaging.GetID() {
			return nil
		}
		if spk.textToSpeechTransformer != nil {
			ctx, span, _ := spk.Tracer().StartSpan(spk.Context(), utils.AssistantSpeakingStage)
			defer span.EndSpan(ctx, utils.AssistantSpeakingStage)
			span.AddAttributes(ctx,
				internal_adapter_telemetry.MessageKV(res.ContextID),
				internal_adapter_telemetry.KV{K: "activity", V: internal_adapter_telemetry.StringValue("speak")},
				internal_adapter_telemetry.KV{K: "script", V: internal_adapter_telemetry.StringValue(res.Text)},
			)
			if err := spk.textToSpeechTransformer.Transform(spk.Context(), res); err != nil {
				spk.logger.Errorf("speak: failed to send flush to text to speech transformer error: %v", err)
			}
		}
		if err := spk.Notify(ctx, &protos.ConversationAssistantMessage{Time: timestamppb.Now(), Id: res.ContextId(), Completed: true, Message: &protos.ConversationAssistantMessage_Text{Text: res.Text}}); err != nil {
			spk.logger.Tracef(ctx, "error while outputting chunk to the user: %w", err)
		}
	default:
	}
	return nil
}

func (talking *GenericRequestor) callDirective(ctx context.Context, vl internal_type.DirectivePacket) error {
	anyArgs, _ := utils.InterfaceMapToAnyMap(vl.Arguments)
	switch vl.Directive {
	case protos.ConversationDirective_END_CONVERSATION:
		if err := talking.Notify(ctx, &protos.AssistantTalkOutput_Directive{Directive: &protos.ConversationDirective{Id: vl.ContextID, Type: vl.Directive, Args: anyArgs, Time: timestamppb.Now()}}); err != nil {
			talking.logger.Errorf("error notifying end conversation action: %v", err)
		}
		return nil
	default:
	}
	return nil
}

/**/
func (talking *GenericRequestor) OnPacket(ctx context.Context, pkts ...internal_type.Packet) error {
	for _, p := range pkts {
		switch vl := p.(type) {
		case internal_type.UserTextPacket:
			// interrupting
			talking.OnPacket(ctx, internal_type.InterruptionPacket{ContextID: vl.ContextID, Source: internal_type.InterruptionSourceWord})

			// add new ID for user text message
			vl.ContextID = talking.messaging.GetID()

			// calling end of speech analyzer
			if err := talking.callEndOfSpeech(ctx, vl); err != nil {
				talking.OnPacket(ctx, internal_type.EndOfSpeechPacket{ContextID: talking.messaging.GetID(), Speech: vl.Text})
			}
			continue

		case internal_type.UserAudioPacket:
			if talking.denoiser != nil && !vl.NoiseReduced {
				vl.NoiseReduced = true
				dnOut, _, err := talking.denoiser.Denoise(ctx, vl.Audio)
				if err != nil {
					talking.logger.Warnf("error while denoising process | will process actual audio byte")
					talking.OnPacket(ctx, vl)
				} else {
					vl.Audio = dnOut
					talking.OnPacket(ctx, vl)
				}
				continue
			}

			if err := talking.callRecording(ctx, vl); err != nil {
				talking.logger.Errorf("recorder error: %v", err)
			}

			if err := talking.callVadProcess(ctx, vl); err != nil {
				talking.logger.Errorf("VAD process error: %v", err)
			}

			if err := talking.callSpeechToText(ctx, vl); err != nil {
				talking.logger.Errorf("speech to text transform error: %v", err)
			}
			continue
		case internal_type.StaticPacket:
			// when static packet is received it means that rapida system has something to speak
			// do not abrupt it just send it to the assembler

			// starting the timer for idle timeout as bot has finished responding
			if talking.messaging.GetInputMode().Text() {
				// stop idle timeout as bot has started responding
				talking.startIdleTimeoutTimer(talking.Context())
			}

			if err := talking.callCreateMessage(ctx, vl); err != nil {
				talking.logger.Errorf("unable to create message from static packet %v", err)
			}

			// sending static packat to executor for any post processing
			if err := talking.messaging.Transition(internal_adapter_request_customizers.LLMGenerating); err != nil {
				talking.logger.Errorf("messaging transition error: %v", err)
			}

			//
			if err := talking.assistantExecutor.Execute(ctx, talking, vl); err != nil {
				talking.logger.Errorf("assistant executor error: %v", err)
			}

			if err := talking.callTextAggregator(ctx, internal_type.LLMResponseDeltaPacket{ContextID: vl.ContextId(), Text: vl.Text}); err != nil {
				talking.logger.Debugf("unable to send static packet to aggregator %v", err)
			}

			if err := talking.messaging.Transition(internal_adapter_request_customizers.LLMGenerated); err != nil {
				talking.logger.Errorf("messaging transition error: %v", err)
			}
			if err := talking.callTextAggregator(ctx, internal_type.LLMResponseDonePacket{ContextID: vl.ContextId()}); err != nil {
				talking.logger.Debugf("unable to send static packet to aggregator %v", err)
			}

			continue
		case internal_type.InterruptionPacket:
			ctx, span, _ := talking.Tracer().StartSpan(talking.Context(), utils.AssistantUtteranceStage)
			defer span.EndSpan(ctx, utils.AssistantUtteranceStage)

			// calling end of speech analyzer
			if err := talking.callEndOfSpeech(ctx, vl); err != nil {
				talking.logger.Errorf("end of speech error: %v", err)
			}
			//
			// recorder interrupted
			if err := talking.callRecording(ctx, vl); err != nil {
				talking.logger.Errorf("recorder error: %v", err)
			}

			// let all the providers know about interruption
			if err := talking.interruptAllProvider(ctx, vl); err != nil {
				talking.logger.Errorf("interrupt all provider error: %v", err)
			}

			switch vl.Source {
			case internal_type.InterruptionSourceWord:
				span.AddAttributes(ctx, internal_telemetry.KV{K: "activity_type", V: internal_telemetry.StringValue("word_interrupt")})
				talking.resetIdleTimeoutTimer(talking.Context())
				//
				if err := talking.messaging.Transition(internal_adapter_request_customizers.Interrupted); err != nil {
					continue
				}

				// notify interruption without waiting
				utils.Go(ctx, func() {
					talking.Notify(ctx, &protos.ConversationInterruption{Type: protos.ConversationInterruption_INTERRUPTION_TYPE_WORD, Time: timestamppb.Now()})
				})
				continue
			default:
				// might be noise at first
				if vl.StartAt < 3 {
					continue
				}
				span.AddAttributes(ctx, internal_telemetry.KV{K: "activity_type", V: internal_telemetry.StringValue("vad_interrupt")})
				if err := talking.messaging.Transition(internal_adapter_request_customizers.Interrupt); err != nil {
					continue
				}

				// notify interruption without waiting
				utils.Go(ctx, func() {
					talking.Notify(ctx, &protos.ConversationInterruption{Type: protos.ConversationInterruption_INTERRUPTION_TYPE_VAD, Time: timestamppb.Now()})
				})
				continue
			}
		case internal_type.SpeechToTextPacket:
			ctx, span, _ := talking.Tracer().StartSpan(talking.Context(), utils.AssistantListeningStage,
				internal_telemetry.KV{
					K: "transcript",
					V: internal_telemetry.StringValue(vl.Script),
				}, internal_telemetry.KV{
					K: "confidence",
					V: internal_telemetry.FloatValue(vl.Confidence),
				}, internal_telemetry.KV{
					K: "isCompleted",
					V: internal_telemetry.BoolValue(!vl.Interim),
				})
			defer span.EndSpan(ctx, utils.AssistantListeningStage)
			// later move the contextID with audio
			vl.ContextID = talking.messaging.GetID()
			//
			if err := talking.callEndOfSpeech(ctx, vl); err != nil {
				if !vl.Interim {
					talking.OnPacket(ctx, internal_type.EndOfSpeechPacket{ContextID: vl.ContextID, Speech: vl.Script})
				}
			}
			continue
		case internal_type.InterimEndOfSpeechPacket:
			talking.Notify(ctx, &protos.ConversationUserMessage{Id: vl.ContextID, Message: &protos.ConversationUserMessage_Text{Text: vl.Speech}, Completed: false, Time: timestamppb.New(time.Now())})
			continue
		case internal_type.EndOfSpeechPacket:
			ctx, span, _ := talking.Tracer().StartSpan(talking.Context(), utils.AssistantUtteranceStage)
			span.EndSpan(ctx,
				utils.AssistantUtteranceStage,
				internal_telemetry.KV{K: "activity_type", V: internal_telemetry.StringValue("SpeechEndActivity")},
				internal_telemetry.KV{K: "speech", V: internal_telemetry.StringValue(vl.Speech)},
			)

			// stop idle timeout as bot has started responding
			talking.stopIdleTimeoutTimer()

			if err := talking.messaging.Transition(internal_adapter_request_customizers.LLMGenerating); err != nil {
				talking.logger.Errorf("messaging transition error: %v", err)
			}

			if err := talking.Notify(ctx,
				&protos.ConversationUserMessage{Id: vl.ContextID, Message: &protos.ConversationUserMessage_Text{Text: vl.Speech}, Completed: true, Time: timestamppb.New(time.Now())}); err != nil {
				talking.logger.Tracef(ctx, "might be returing processing the duplicate message so cut it out.")
				continue
			}
			utils.Go(ctx, func() {
				if err := talking.onCreateMessage(ctx, internal_type.UserTextPacket{ContextID: vl.ContextID, Text: vl.Speech}); err != nil {
					talking.logger.Errorf("Error in onCreateMessage: %v", err)
				}
			})

			//
			if err := talking.assistantExecutor.Execute(ctx, talking, internal_type.UserTextPacket{ContextID: vl.ContextID, Text: vl.Speech}); err != nil {
				talking.logger.Errorf("assistant executor error: %v", err)
				talking.OnError(ctx)
				continue
			}
		case internal_type.LLMResponseDeltaPacket:

			// might be stale packet
			if vl.ContextID != talking.messaging.GetID() {
				continue
			}

			if err := talking.messaging.Transition(internal_adapter_request_customizers.LLMGenerating); err != nil {
				talking.logger.Errorf("messaging transition error: %v", err)
			}
			// sending to aggregator for assembling sentences
			if err := talking.callTextAggregator(ctx, vl); err != nil {
				talking.logger.Errorf("sentence aggregator error: %v, calling speak directly", err)
				if err := talking.callSpeaking(ctx, vl); err != nil {
					talking.logger.Errorf("speaking error: %v", err)
				}
			}
			// end of speech analyzer in case histoyrical data is to be used

		case internal_type.LLMResponseDonePacket:

			// might be stale packet
			if vl.ContextID != talking.messaging.GetID() {
				continue
			}

			// starting the timer for idle timeout as bot has finished responding
			if talking.messaging.GetInputMode().Text() {
				// stop idle timeout as bot has started responding
				talking.startIdleTimeoutTimer(talking.Context())
			}
			//
			if err := talking.messaging.Transition(internal_adapter_request_customizers.LLMGenerated); err != nil {
				talking.logger.Errorf("messaging transition error: %v", err)
			}
			if err := talking.callCreateMessage(ctx, vl); err != nil {
				talking.logger.Errorf("error creating message: %v", err)
			}

			if err := talking.callTextAggregator(ctx, vl); err != nil {
				talking.logger.Errorf("sentence aggregator error: %v calling speak directly", err)
				if err := talking.callSpeaking(ctx, vl); err != nil {
					talking.logger.Errorf("speaking error: %v", err)
				}
			}

			continue

		case internal_type.DirectivePacket:
			talking.callDirective(ctx, vl)
			continue
		case internal_type.MetricPacket:
			// metrics update for the message
			// later this can be used at each stage to calculate various metrics
			if len(vl.Metrics) > 0 {
				if err := talking.onMessageMetric(talking.Context(), vl.ContextID, vl.Metrics); err != nil {
					talking.logger.Errorf("Error in OnUpdateMessage: %v", err)
				}
			}
		case internal_type.TextToSpeechEndPacket:
			// might be stale packet
			if vl.ContextID != talking.messaging.GetID() {
				continue
			}
			if err := talking.Notify(talking.Context(), &protos.ConversationAssistantMessage{Time: timestamppb.Now(), Id: vl.ContextID, Completed: true}); err != nil {
				talking.logger.Tracef(talking.ctx, "error while outputing chunk to the user: %w", err)
			}

			continue
		case internal_type.TextToSpeechAudioPacket:

			// resetting idle timer as bot has sponken
			// starting the timer for idle timeout as bot has finished responding
			if talking.messaging.GetInputMode().Audio() {
				// stop idle timeout as bot has started responding
				talking.startIdleTimeoutTimer(talking.Context())
			}

			// might be stale packet
			if vl.ContextID != talking.messaging.GetID() {
				continue
			}

			// notify the user about audio chunk
			if err := talking.Notify(talking.Context(), &protos.ConversationAssistantMessage{Time: timestamppb.Now(), Id: vl.ContextID, Message: &protos.ConversationAssistantMessage_Audio{Audio: vl.AudioChunk}, Completed: false}); err != nil {
				talking.logger.Tracef(talking.ctx, "error while outputing chunk to the user: %w", err)
			}

			// for recording puposes
			if err := talking.callRecording(ctx, vl); err != nil {
				talking.logger.Errorf("recorder error: %v", err)
			}
			continue
		default:
			talking.logger.Warnf("unknown packet type received in OnGeneration %T", vl)
		}
	}

	return nil
}
