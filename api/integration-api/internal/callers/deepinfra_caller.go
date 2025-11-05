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

type deepInfrAICaller struct {
	AICaller
	logger commons.Logger
	apiUrl string
}

// CallWithPayload implements Caller.
func (aiC *deepInfrAICaller) CallWithPayload(ctx context.Context, endpoint string, method string, headers map[string]string, payload io.Reader) (*string, error) {
	panic("unimplemented")
}

type DeepInfraError struct {
	Detail struct {
		Error string `json:"error"`
	} `json:"detail"`
}

func (e DeepInfraError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "undefined error"
	}
	return string(b)
}

func NewDeepInfrAICaller(logger commons.Logger, apiUrl string) Caller {
	return &deepInfrAICaller{
		AICaller: AICaller{
			logger: logger,
		},
		logger: logger,
		apiUrl: apiUrl,
	}
}

func (aiC *deepInfrAICaller) Call(ctx context.Context, endpoint, method string, headers map[string]string, payload map[string]interface{}) (*string, error) {
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

func (aiC *deepInfrAICaller) do(req *http.Request) (*string, error) {

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

	// It gives different responses for different error, keep it as string for now
	var apiErr DeepInfraError
	if err := Unmarshal(resp, &apiErr); err != nil {
		return nil, err
	}

	return nil, apiErr
}

func (aiC *deepInfrAICaller) Endpoint(url string) string {
	return fmt.Sprintf("%s%s", aiC.apiUrl, url)
}

func (aiC *deepInfrAICaller) Header(req *http.Request, headers map[string]string) {
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}
