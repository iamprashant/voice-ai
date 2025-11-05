package internal_bedrock_callers

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	internal_callers "github.com/rapidaai/internal/callers"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	integration_api "github.com/rapidaai/protos"
)

type verifyCredentialCaller struct {
	Bedrock
}

func NewVerifyCredentialCaller(logger commons.Logger, credential *integration_api.Credential) internal_callers.Verifier {
	return &verifyCredentialCaller{
		Bedrock: bedrock(logger, credential),
	}
}

func (stc *verifyCredentialCaller) CredentialVerifier(
	ctx context.Context,
	options *internal_callers.CredentialVerifierOptions) (*string, error) {
	cfg, err := stc.Cfg()
	if err != nil {
		return nil, err
	}

	// single minute timeout and cancellable by the client as context will get cancel
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Call the GetCallerIdentity API
	_, err = sts.NewFromConfig(*cfg).GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return utils.Ptr("valid"), nil
}
