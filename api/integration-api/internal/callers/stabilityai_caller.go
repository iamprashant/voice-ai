package internal_callers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rapidaai/pkg/commons"
)

type stabilityAiCaller struct {
	AICaller
	logger commons.Logger
	apiUrl string
}

// CallWithPayload implements Caller.
func (aiC *stabilityAiCaller) CallWithPayload(ctx context.Context, endpoint string, method string, headers map[string]string, payload io.Reader) (*string, error) {
	panic("unimplemented")
}

func NewStabilityAiCaller(logger commons.Logger, apiUrl string) Caller {
	return &stabilityAiCaller{
		AICaller: AICaller{
			logger: logger,
		},
		logger: logger,
		apiUrl: apiUrl,
	}
}

type StabilityAiError struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e StabilityAiError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "undefined error"
	}
	return string(b)
}

func (aiC *stabilityAiCaller) Call(ctx context.Context, endpoint, method string, headers map[string]string, payload map[string]interface{}) (*string, error) {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, aiC.Endpoint(endpoint), bytes.NewBuffer(encodedPayload))
	if err != nil {
		return nil, err
	}
	aiC.Header(req, headers)
	return aiC.do(req)

}

func (aiC *stabilityAiCaller) do(req *http.Request) (*string, error) {

	resp, err := aiC.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Check for valid status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			aiC.logger.Errorf("unable to read all the response from byte %v", err)
		}
		bodyString := string(bodyBytes)
		return &bodyString, nil
	}

	// If we have not-success HTTP status code, unmarshal to APIError
	var apiErr StabilityAiError
	if err := Unmarshal(resp, &apiErr); err != nil {
		return nil, err
	}
	if apiErr.StatusCode == 0 {
		// Overwrite apiErr status code if it's zero
		apiErr.StatusCode = resp.StatusCode
	}
	return nil, apiErr
}

func (aiC *stabilityAiCaller) Endpoint(url string) string {
	return fmt.Sprintf("%s%s", aiC.apiUrl, url)
}

func (aiC *stabilityAiCaller) Header(req *http.Request, headers map[string]string) {
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}
