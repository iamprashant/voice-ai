package internal_bedrock_callers

import (
	"context"
	"errors"

	internal_callers "github.com/rapidaai/internal/callers"
	"github.com/rapidaai/pkg/commons"
	integration_api "github.com/rapidaai/protos"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/smithy-go/logging"
)

type Bedrock struct {
	logger     commons.Logger
	credential internal_callers.CredentialResolver
}

// {"name": "access_key_id", "type": "string", "label": "Access Key ID"},
// {"name": "secret_access_key", "type": "string", "label": "Secret Access Key"},
// {"name": "region", "type": "string", "label": "Region"}

var (
	REGION        = "region"
	ACCESS_KEY_ID = "access_key_id"
	SECRET_KEY    = "secret_access_key"
)

func bedrock(logger commons.Logger, credential *integration_api.Credential) Bedrock {
	return Bedrock{
		logger: logger,
		credential: func() map[string]interface{} {
			return credential.GetValue().AsMap()
		},
	}
}

func (br *Bedrock) Logf(classification logging.Classification, format string, v ...interface{}) {
	br.logger.Debugf(format, v)
}

func (br *Bedrock) Cfg() (*aws.Config, error) {
	crds := br.credential()

	region, ok := crds[REGION]
	if !ok {
		br.logger.Errorf("Unable to get client for bedrock without region")
		return nil, errors.New("unable to resolve the credential for aws")
	}

	accessKeyId, ok := crds[ACCESS_KEY_ID]
	if !ok {
		br.logger.Errorf("Unable to get client for bedrock without accessKeyId")
		return nil, errors.New("unable to resolve the credential for aws")
	}

	secretKey, ok := crds[SECRET_KEY]
	if !ok {
		br.logger.Errorf("Unable to get client for bedrock without secretKey")
		return nil, errors.New("unable to resolve the credential for aws")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region.(string)),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyId.(string), secretKey.(string), ""),
		),
		config.WithLogger(br),
	)
	if err != nil {
		br.logger.Errorf("Unable to get client for bedroxk")
		return nil, errors.New("unable to resolve the credential")
	}
	return &cfg, err
}

func (br *Bedrock) GetClient() (*bedrockruntime.Client, error) {
	cfg, err := br.Cfg()
	if err != nil {
		br.logger.Errorf("Unable to get aws config")
		return nil, errors.New("unable to resolve the credential")
	}
	return bedrockruntime.NewFromConfig(*cfg), nil

}

type AWSCompletionResponse struct {
	Completions []*AWSCompletion `json:"completions"`
}
type AWSCompletion struct {
	Data *AWSCompletionData `json:"data"`
}
type AWSCompletionData struct {
	Text string `json:"text"`
}

type AWSConverseOutput struct {
	Message struct {
		Role    string `json:"role"`
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

func (AWSConverseOutput) isConverseOutput() {}

type AWSConverseMetrics struct {
	LatencyMs int `json:"latencyMs"`
}
type AWSConverseUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}
