// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package channel_telephony

import (
	"bufio"
	"context"
	"errors"
	"net"

	"github.com/gorilla/websocket"
	"github.com/rapidaai/api/assistant-api/config"
	internal_asterisk_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/asterisk"
	internal_asterisk_audiosocket "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/asterisk/audiosocket"
	internal_asterisk_websocket "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/asterisk/websocket"
	internal_exotel_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/exotel"
	internal_sip_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/sip"
	internal_twilio_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/twilio"
	internal_vonage_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/vonage"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	sip_infra "github.com/rapidaai/api/assistant-api/sip/infra"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type Telephony string

const (
	Twilio   Telephony = "twilio"
	Exotel   Telephony = "exotel"
	Vonage   Telephony = "vonage"
	Asterisk Telephony = "asterisk"
	SIP      Telephony = "sip"
)

func (at Telephony) String() string {
	return string(at)
}

func GetTelephony(at Telephony, cfg *config.AssistantConfig, logger commons.Logger, sipServer *sip_infra.Server) (internal_type.Telephony, error) {
	switch at {
	case Twilio:
		return internal_twilio_telephony.NewTwilioTelephony(cfg, logger)
	case Exotel:
		return internal_exotel_telephony.NewExotelTelephony(cfg, logger)
	case Vonage:
		return internal_vonage_telephony.NewVonageTelephony(cfg, logger)
	case Asterisk:
		return internal_asterisk_telephony.NewAsteriskTelephony(cfg, logger)
	case SIP:
		return internal_sip_telephony.NewSIPTelephony(cfg, logger, sipServer)
	default:
		return nil, errors.New("illegal telephony provider")
	}
}

// streamer

// Streamer creates a new WebSocket streamer for Asterisk chan_websocket
func (at Telephony) Streamer(
	ctx context.Context,
	logger commons.Logger,
	connection *websocket.Conn,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
) (internal_type.TelephonyStreamer, error) {
	switch at {
	case Twilio:
		return internal_twilio_telephony.NewTwilioWebsocketStreamer(logger, connection, assistant, conversation, vlt), nil
	case Exotel:
		return internal_exotel_telephony.NewExotelWebsocketStreamer(logger, connection, assistant, conversation, vlt), nil
	case Vonage:
		return internal_vonage_telephony.NewVonageWebsocketStreamer(logger, connection, assistant, conversation, vlt), nil
	case Asterisk:
		return internal_asterisk_websocket.NewAsteriskWebsocketStreamer(logger, connection, assistant, conversation, vlt), nil
	default:
		return nil, errors.New("illegal telephony provider")
	}

}

func (at Telephony) AudioSocketStreamer(
	logger commons.Logger,
	conn net.Conn,
	reader *bufio.Reader,
	writer *bufio.Writer,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
) (internal_type.TelephonyStreamer, error) {
	switch at {
	case Asterisk:
		return internal_asterisk_audiosocket.NewStreamer(logger, conn, reader, writer, assistant, conversation, vlt)
	}
	return nil, errors.New("illegal telephony provider")
}

func (at Telephony) SipStreamer(
	ctx context.Context,
	logger commons.Logger,
	session *sip_infra.Session,
	config *sip_infra.Config,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
) (internal_type.TelephonyStreamer, error) {
	switch at {
	case SIP:
		return internal_sip_telephony.NewInboundStreamer(ctx, config, logger, session, assistant, conversation)
	}
	return nil, errors.New("illegal telephony provider")
}
