// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_exotel_telephony

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rapidaai/api/assistant-api/config"
	internal_exotel "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/exotel/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

const exotelProvider = "exotel"

type exotelTelephony struct {
	logger commons.Logger
	appCfg *config.AssistantConfig
}

func NewExotelTelephony(config *config.AssistantConfig, logger commons.Logger) (internal_type.Telephony, error) {
	return &exotelTelephony{
		logger: logger,
		appCfg: config,
	}, nil
}

func (tpc *exotelTelephony) CatchAllStatusCallback(ctx *gin.Context) (*internal_type.StatusInfo, error) {
	return nil, nil
}

// StatusCallback implements [Telephony].
func (tpc *exotelTelephony) StatusCallback(c *gin.Context, auth types.SimplePrinciple, assistantId uint64, assistantConversationId uint64) (*internal_type.StatusInfo, error) {
	form, err := c.MultipartForm()
	if err != nil {
		tpc.logger.Errorf("failed to parse multipart form-data with error %+v", err)
		return nil, fmt.Errorf("failed to parse multipart form-data")
	}

	eventDetails := make(map[string]interface{})
	for key, values := range form.Value {
		if len(values) > 0 {
			eventDetails[key] = values[0]
		} else {
			eventDetails[key] = nil
		}
	}
	event := fmt.Sprintf("%v", eventDetails["Status"])
	return &internal_type.StatusInfo{Event: event, Payload: eventDetails}, nil
}
func (tpc *exotelTelephony) ClientUrl(vaultCredential *protos.VaultCredential, opts utils.Option) (*string, error) {
	accountSid, ok := vaultCredential.GetValue().AsMap()["account_sid"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config accountSid is not found")
	}
	clientId, ok := vaultCredential.GetValue().AsMap()["client_id"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config client_id not found")
	}
	authToken, ok := vaultCredential.GetValue().AsMap()["client_secret"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config")
	}
	return utils.Ptr(fmt.Sprintf("https://%s:%s@api.exotel.com/v1/Accounts/%s/Calls/connect.json",
		clientId.(string), authToken.(string), accountSid.(string))), nil

}

func (tpc *exotelTelephony) AppUrl(
	vaultCredential *protos.VaultCredential,
	opts utils.Option) (*string, error) {
	accountSid, ok := vaultCredential.GetValue().AsMap()["account_sid"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config accountSid is not found")
	}
	app_id, err := opts.GetString("app_id")
	if err != nil {
		return nil, fmt.Errorf("illegal app_id option is not found")
	}
	return utils.Ptr(fmt.Sprintf("http://my.exotel.com/%s/exoml/start_voice/%s", accountSid.(string), app_id)), nil

}

func (tpc *exotelTelephony) OutboundCall(
	auth types.SimplePrinciple,
	// customer number
	toPhone string,
	// exo number
	fromPhone string,
	assistantId, assistantConversationId uint64,
	vaultCredential *protos.VaultCredential,
	opts utils.Option) (*internal_type.CallInfo, error) {
	info := &internal_type.CallInfo{Provider: exotelProvider}

	clientUrl, err := tpc.ClientUrl(vaultCredential, opts)
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("Failed to build url, check credentials: %s", err.Error())
		return info, err
	}

	appUrl, err := tpc.AppUrl(vaultCredential, opts)
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("Failed to build app url: %s", err.Error())
		return info, err
	}

	contextID, _ := opts.GetString("rapida.context_id")

	formData := url.Values{}
	formData.Set("From", toPhone)
	formData.Set("CallerId", fromPhone)
	formData.Set("To", fromPhone)
	formData.Set("Url", *appUrl)
	formData.Set("StatusCallback", fmt.Sprintf("https://%s/%s", tpc.appCfg.PublicAssistantHost, internal_type.GetContextEventPath(exotelProvider, contextID)))
	// for exotel there is no way to set dynamic path so pass it as custom filed
	formData.Set("CustomField", internal_type.GetContextAnswerPath(exotelProvider, contextID))

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", *clientUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("request creation error: %s", err.Error())
		return info, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("API error: %s", err.Error())
		return info, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("failed to read response: %s", err.Error())
		return info, err
	}
	if resp.StatusCode != http.StatusOK {
		tpc.logger.Errorf("Unexpected HTTP Status: %d, Response Body: %s\n", resp.StatusCode, string(bodyBytes))
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
		info.StatusInfo = internal_type.StatusInfo{Event: "Failed", Payload: string(bodyBytes)}
		return info, fmt.Errorf("status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var jsonResponse internal_exotel.MakeCallResponse
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("failed to decode response: %s", err.Error())
		info.StatusInfo = internal_type.StatusInfo{Event: jsonResponse.Call.Status, Payload: "Failed to decode response"}
		return info, err
	}

	info.ChannelUUID = jsonResponse.Call.Sid
	info.Status = "SUCCESS"
	info.StatusInfo = internal_type.StatusInfo{Event: jsonResponse.Call.Status, Payload: jsonResponse}
	return info, nil
}

func (tpc *exotelTelephony) InboundCall(c *gin.Context, auth types.SimplePrinciple, assistantId uint64, clientNumber string, assistantConversationId uint64) error {
	contextID, _ := c.Get("contextId")
	ctxID := fmt.Sprintf("%v", contextID)

	response := map[string]string{
		"url": fmt.Sprintf("wss://%s/%s",
			tpc.appCfg.PublicAssistantHost,
			internal_type.GetContextAnswerPath("exotel", ctxID)),
	}
	c.JSON(http.StatusOK, response)
	return nil
}

func (tpc *exotelTelephony) ReceiveCall(c *gin.Context) (*internal_type.CallInfo, error) {
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	socketUrl, ok := queryParams["CustomField"]
	if ok {
		response := map[string]string{"url": fmt.Sprintf("wss://%s/%s", tpc.appCfg.PublicAssistantHost, socketUrl)}
		c.JSON(http.StatusOK, response)
		return nil, fmt.Errorf("outbound call triggered")
	}

	clientNumber, ok := queryParams["CallFrom"]
	if !ok || clientNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid caller"})
		return nil, fmt.Errorf("missing or empty 'from' query parameter")
	}

	info := &internal_type.CallInfo{
		CallerNumber: clientNumber,
		Provider:     exotelProvider,
		Status:       "SUCCESS",
		StatusInfo:   internal_type.StatusInfo{Event: "webhook", Payload: queryParams},
	}
	if v, ok := queryParams["CallSid"]; ok && v != "" {
		info.ChannelUUID = v
	}
	return info, nil
}
