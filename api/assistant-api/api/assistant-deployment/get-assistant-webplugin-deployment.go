package assistant_deployment_api

import (
	"context"
	"errors"

	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	assistant_api "github.com/rapidaai/protos"
	lexatic_backend "github.com/rapidaai/protos"
)

// GetAssistantWebpluginDeployment implements lexatic_backend.AssistantDeploymentServiceServer.
func (deploymentApi *assistantDeploymentGrpcApi) GetAssistantWebpluginDeployment(ctx context.Context, getter *lexatic_backend.GetAssistantDeploymentRequest) (*lexatic_backend.GetAssistantWebpluginDeploymentResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || iAuth.GetCurrentProjectId() == nil {
		deploymentApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[assistant_api.GetAssistantWebpluginDeploymentResponse](
			errors.New("unauthenticated request for create assistant whatsapp deployment"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	webpluginDeployment, err := deploymentApi.deploymentService.GetAssistantWebpluginDeployment(ctx, iAuth, getter.GetAssistantId())
	if err != nil {
		return utils.Error[lexatic_backend.GetAssistantWebpluginDeploymentResponse](err, "Unable to get deployment, please try again later.")
	}
	out := &lexatic_backend.AssistantWebpluginDeployment{}
	err = utils.Cast(webpluginDeployment, out)
	if err != nil {
		deploymentApi.logger.Warnf("unable to cast the web plugin deployment model to the response object")
	}
	return &lexatic_backend.GetAssistantWebpluginDeploymentResponse{
		Data:    out,
		Success: true,
		Code:    200,
	}, nil
}
