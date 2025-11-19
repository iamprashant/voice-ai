package assistant_deployment_api

import (
	"context"
	"errors"

	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	assistant_api "github.com/rapidaai/protos"
)

// CreateAssistantWebpluginDeployment implements assistant_api.AssistantDeploymentServiceServer.
func (deploymentApi *assistantDeploymentApi) CreateAssistantWebpluginDeployment(ctx context.Context, deployment *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantWebpluginDeploymentResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || iAuth.GetCurrentProjectId() == nil {
		deploymentApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[assistant_api.GetAssistantWebpluginDeploymentResponse](
			errors.New("unauthenticated request for create assistant web plugin deployment"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}

	if deployment.GetPlugin() == nil {
		return utils.Error[assistant_api.GetAssistantWebpluginDeploymentResponse](
			errors.New("illegal parameters attached to deployment"),
			"Please check and provide valid deployment request for webplugin.",
		)
	}
	iconPath := "https://cdn-01.rapida.ai/partners/rapida.png"
	switch icon := deployment.GetPlugin().GetIcon().(type) {
	case *assistant_api.AssistantWebpluginDeployment_Raw:
		so := deploymentApi.storage.Store(ctx,
			icon.Raw.GetName(),
			icon.Raw.GetContent(),
		)
		if so.Error != nil {
			deploymentApi.logger.Errorf("error while uploading the image to cdn %+v", so.Error)
			break
		}
		iconPath = so.CompletePath
	case *assistant_api.AssistantWebpluginDeployment_Url:
		iconPath = icon.Url
		break
	}

	wpDeployment, err := deploymentApi.deploymentService.CreateWebPluginDeployment(ctx,
		iAuth, deployment.GetPlugin().GetAssistantId(),
		deployment.GetPlugin().GetName(),
		iconPath,
		deployment.GetPlugin().Greeting,
		deployment.GetPlugin().Mistake,
		&deployment.GetPlugin().IdealTimeout,
		&deployment.GetPlugin().IdealTimeoutMessage,
		&deployment.GetPlugin().MaxSessionDuration,
		deployment.GetPlugin().GetSuggestion(),
		deployment.GetPlugin().GetHelpCenterEnabled(),
		deployment.GetPlugin().GetProductCatalogEnabled(),
		deployment.GetPlugin().GetArticleCatalogEnabled(),
		deployment.GetPlugin().GetInputAudio(),
		deployment.GetPlugin().GetOutputAudio(),
	)

	if err != nil {
		return utils.Error[assistant_api.GetAssistantWebpluginDeploymentResponse](
			errors.New("unauthenticated request for create assistant webplugin deployment"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	return utils.Success[assistant_api.GetAssistantWebpluginDeploymentResponse](wpDeployment)
}
