package workflow_client

import (
	"context"
	"time"

	"github.com/rapidaai/config"
	clients "github.com/rapidaai/pkg/clients"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/types"
	assistant_api "github.com/rapidaai/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AssistantServiceClient interface {
	GetAllAssistant(c context.Context, auth types.SimplePrinciple, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.Assistant, error)

	DeleteAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.DeleteAssistantRequest) (*assistant_api.GetAssistantResponse, error)
	GetAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantRequest) (*assistant_api.GetAssistantResponse, error)
	CreateAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantRequest) (*assistant_api.GetAssistantResponse, error)

	GetAllAssistantProvider(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.GetAllAssistantProviderResponse_AssistantProvider, error)
	UpdateAssistantVersion(c context.Context, auth types.SimplePrinciple, assistantId, assistantProviderModelId uint64) (*assistant_api.GetAssistantResponse, error)
	CreateAssistantProvider(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantProviderRequest) (*assistant_api.GetAssistantProviderResponse, error)

	//
	GetAllMessage(c context.Context, auth types.SimplePrinciple,
		criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate,
		order *assistant_api.Ordering, selectors []*assistant_api.FieldSelector) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error)
	GetAllAssistantMessage(c context.Context, auth types.SimplePrinciple, assistantId uint64,
		criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate,
		order *assistant_api.Ordering, selectors []*assistant_api.FieldSelector) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error)
	GetAllAssistantConversation(ctx context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, order *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantConversation, error)
	GetAllConversationMessage(ctx context.Context, auth types.SimplePrinciple, assistantId, assistantConversationId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, order *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error)
	GetAssistantConversation(
		c context.Context,
		auth types.SimplePrinciple,
		assistantRequest *assistant_api.GetAssistantConversationRequest) (*assistant_api.GetAssistantConversationResponse, error)

	CreateAssistantTag(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantTagRequest) (*assistant_api.GetAssistantResponse, error)
	UpdateAssistantDetail(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.UpdateAssistantDetailRequest) (*assistant_api.GetAssistantResponse, error)

	// deployment
	CreateAssistantApiDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantApiDeploymentResponse, error)
	CreateAssistantPhoneDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantPhoneDeploymentResponse, error)
	CreateAssistantWhatsappDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantWhatsappDeploymentResponse, error)
	CreateAssistantWebpluginDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantWebpluginDeploymentResponse, error)
	CreateAssistantDebuggerDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantDebuggerDeploymentResponse, error)

	GetAssistantApiDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantApiDeploymentResponse, error)
	GetAssistantPhoneDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantPhoneDeploymentResponse, error)
	GetAssistantWhatsappDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantWhatsappDeploymentResponse, error)
	GetAssistantWebpluginDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantWebpluginDeploymentResponse, error)
	GetAssistantDebuggerDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantDebuggerDeploymentResponse, error)

	//
	GetAssistantWebhookLog(ctx context.Context, auth types.SimplePrinciple, req *assistant_api.GetAssistantWebhookLogRequest) (*assistant_api.GetAssistantWebhookLogResponse, error)
	GetAllAssistantWebhookLog(ctx context.Context, auth types.SimplePrinciple, projectId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, ordering *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantWebhookLog, error)

	//
	GetAllAssistantWebhook(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantWebhook, error)
	GetAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error)
	CreateAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error)
	UpdateAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error)
	DeleteAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error)

	//
	GetAllAssistantAnalysis(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantAnalysis, error)
	GetAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error)
	CreateAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error)
	UpdateAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error)
	DeleteAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error)

	//
	GetAllAssistantTool(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantTool, error)
	GetAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error)
	CreateAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error)
	UpdateAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error)
	DeleteAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error)

	//
	GetAllAssistantKnowledge(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantKnowledge, error)
	GetAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error)
	CreateAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error)
	UpdateAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error)
	DeleteAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error)

	GetAssistantToolLog(ctx context.Context, auth types.SimplePrinciple, in *assistant_api.GetAssistantToolLogRequest) (*assistant_api.GetAssistantToolLogResponse, error)
	GetAllAssistantToolLog(ctx context.Context, auth types.SimplePrinciple, in *assistant_api.GetAllAssistantToolLogRequest) (*assistant_api.GetAllAssistantToolLogResponse, error)
	GetAllAssistantTelemetry(ctx context.Context, auth types.SimplePrinciple, in *assistant_api.GetAllAssistantTelemetryRequest) (*assistant_api.GetAllAssistantTelemetryResponse, error)
}

type assistantServiceClient struct {
	clients.InternalClient
	cfg                       *config.AppConfig
	logger                    commons.Logger
	assistantClient           assistant_api.AssistantServiceClient
	assistantDeploymentClient assistant_api.AssistantDeploymentServiceClient
}

// NewAssistantServiceClientGRPC creates a new instance of AssistantServiceClient using gRPC.
// It establishes a connection to the assistant service using the provided configuration, logger, and Redis connector.
//
// Parameters:
// - config: The application configuration containing the workflow host details.
// - logger: A Logger instance for logging messages.
// - redis: A RedisConnector instance for connecting to Redis.
//
// Returns:
// - An instance of AssistantServiceClient, or nil if an error occurs during connection establishment.
func NewAssistantServiceClientGRPC(config *config.AppConfig, logger commons.Logger, redis connectors.RedisConnector) AssistantServiceClient {
	logger.Debugf("conntecting to assistant client with %s", config.WorkflowHost)
	conn, err := grpc.NewClient(config.WorkflowHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Unable to create connection %v", err)
	}
	return &assistantServiceClient{
		InternalClient:            clients.NewInternalClient(config, logger, redis),
		cfg:                       config,
		logger:                    logger,
		assistantClient:           assistant_api.NewAssistantServiceClient(conn),
		assistantDeploymentClient: assistant_api.NewAssistantDeploymentServiceClient(conn),
	}
}

func (client *assistantServiceClient) GetAllAssistant(c context.Context, auth types.SimplePrinciple, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.Assistant, error) {
	client.logger.Debugf("get all assistant request")
	res, err := client.assistantClient.GetAllAssistant(client.WithAuth(c, auth), &assistant_api.GetAllAssistantRequest{
		Paginate:  paginate,
		Criterias: criterias,
	})
	if err != nil {
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) DeleteAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.DeleteAssistantRequest) (*assistant_api.GetAssistantResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.DeleteAssistant(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.DeleteAssistant", time.Since(start))
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to delete assistant %v", err)
	}
	return res, nil
}

func (client *assistantServiceClient) GetAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantRequest) (*assistant_api.GetAssistantResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistant(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAssistant", time.Since(start))
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get assistant %v", err)
	}
	return res, nil
}

func (client *assistantServiceClient) CreateAssistant(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantRequest) (*assistant_api.GetAssistantResponse, error) {
	res, err := client.assistantClient.CreateAssistant(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Errorf("error while calling CreateAssistant %v", err)
		return nil, err
	}
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantProvider(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.GetAllAssistantProviderResponse_AssistantProvider, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantProvider(client.WithAuth(c, auth), &assistant_api.GetAllAssistantProviderRequest{
		Criterias:   criterias,
		Paginate:    paginate,
		AssistantId: assistantId,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllAssistantProvider", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) UpdateAssistantVersion(c context.Context, auth types.SimplePrinciple, assistantId, assistantProviderModelId uint64) (*assistant_api.GetAssistantResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantVersion(client.WithAuth(c, auth), &assistant_api.UpdateAssistantVersionRequest{
		AssistantId:         assistantId,
		AssistantProviderId: assistantProviderModelId,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.UpdateAssistantVersion", time.Since(start))
		client.logger.Errorf("error while calling to UpdateAssistantVersion %v", err)
		return nil, err
	}
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantProvider(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantProviderRequest) (*assistant_api.GetAssistantProviderResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantProvider(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.CreateAssistantProvider", time.Since(start))
		client.logger.Errorf("error while calling to CreateAssistantProvider %v", err)
		return nil, err
	}
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantTag(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantTagRequest) (*assistant_api.GetAssistantResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantTag(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.UpdateAssistantDetail", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantTag %v", err)
		return nil, err
	}
	return res, nil
}

func (client *assistantServiceClient) UpdateAssistantDetail(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.UpdateAssistantDetailRequest) (*assistant_api.GetAssistantResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantDetail(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.UpdateAssistantDetail", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantTag %v", err)
		return nil, err
	}
	return res, nil
}

func (client *assistantServiceClient) GetAllMessage(ctx context.Context,
	auth types.SimplePrinciple,
	criterias []*assistant_api.Criteria,
	paginate *assistant_api.Paginate,
	order *assistant_api.Ordering,
	fieldSelector []*assistant_api.FieldSelector,
) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllMessage(client.WithAuth(ctx, auth), &assistant_api.GetAllMessageRequest{
		Paginate:  paginate,
		Criterias: criterias,
		Order:     order,
		Selectors: fieldSelector,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllMessage", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllMessage", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllMessage", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAllAssistantMessage(ctx context.Context,
	auth types.SimplePrinciple,
	assistantId uint64, criterias []*assistant_api.Criteria,
	paginate *assistant_api.Paginate,
	order *assistant_api.Ordering,
	fieldSelector []*assistant_api.FieldSelector,
) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantMessage(client.WithAuth(ctx, auth), &assistant_api.GetAllAssistantMessageRequest{
		AssistantId: assistantId,
		Paginate:    paginate,
		Criterias:   criterias,
		Order:       order,
		Selectors:   fieldSelector,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllAssistantMessage", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllAssistantMessage", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAllAssistantConversation(ctx context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, order *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantConversation, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantConversation(client.WithAuth(ctx, auth), &assistant_api.GetAllAssistantConversationRequest{
		AssistantId: assistantId,
		Paginate:    paginate,
		Criterias:   criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllAssistantConversation", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllAssistantConversation", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAllConversationMessage(ctx context.Context, auth types.SimplePrinciple, assistantId, assistantConversationId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, order *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantConversationMessage, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllConversationMessage(client.WithAuth(ctx, auth), &assistant_api.GetAllConversationMessageRequest{
		AssistantConversationId: assistantConversationId,
		AssistantId:             assistantId,
		Paginate:                paginate,
		Criterias:               criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllConversationMessage", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAllConversationMessage", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) CreateAssistantApiDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantApiDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.CreateAssistantApiDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantApiDeployment", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantApiDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantApiDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) CreateAssistantPhoneDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantPhoneDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.CreateAssistantPhoneDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.CreateAssistantPhoneDeployment", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantPhoneDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantServiceClient.CreateAssistantPhoneDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) CreateAssistantWhatsappDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantWhatsappDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.CreateAssistantWhatsappDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantWhatsappDeployment", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantWhatsappDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantWhatsappDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) CreateAssistantWebpluginDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantWebpluginDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.CreateAssistantWebpluginDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantWebpluginDeployment", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantWebpluginDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantWebpluginDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) CreateAssistantDebuggerDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.CreateAssistantDeploymentRequest) (*assistant_api.GetAssistantDebuggerDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.CreateAssistantDebuggerDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantDebuggerDeployment", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantDebuggerDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantDebuggerDeployment", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAssistantApiDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantApiDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.GetAssistantApiDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantApiDeployment", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantApiDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.CreateAssistantDebuggerDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) GetAssistantPhoneDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantPhoneDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.GetAssistantPhoneDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantPhoneDeployment", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantPhoneDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantPhoneDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) GetAssistantWhatsappDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantWhatsappDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.GetAssistantWhatsappDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantWhatsappDeployment", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantWhatsappDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantWhatsappDeployment", time.Since(start))
	return res, nil
}
func (client *assistantServiceClient) GetAssistantWebpluginDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantWebpluginDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.GetAssistantWebpluginDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantWebpluginDeployment", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantWebpluginDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantWebpluginDeployment", time.Since(start))
	client.logger.Debugf("report %+v", res.Data)
	return res, nil
}
func (client *assistantServiceClient) GetAssistantDebuggerDeployment(c context.Context, auth types.SimplePrinciple, assistantRequest *assistant_api.GetAssistantDeploymentRequest) (*assistant_api.GetAssistantDebuggerDeploymentResponse, error) {
	start := time.Now()
	res, err := client.assistantDeploymentClient.GetAssistantDebuggerDeployment(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantDebuggerDeployment", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantDebuggerDeployment %v", err)
		return nil, err
	}
	client.logger.Benchmark("Benchmarking: assistantDeploymentClient.GetAssistantDebuggerDeployment", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantWebhook(ctx context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantWebhook, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantWebhook(client.WithAuth(ctx, auth), &assistant_api.GetAllAssistantWebhookRequest{
		AssistantId: assistantId,
		Paginate:    paginate,
		Criterias:   criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantWebhook", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantWebhook", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAllAssistantWebhookLog(ctx context.Context, auth types.SimplePrinciple,
	projectId uint64,
	criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate, ordering *assistant_api.Ordering) (*assistant_api.Paginated, []*assistant_api.AssistantWebhookLog, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantWebhookLog(client.WithAuth(ctx, auth), &assistant_api.GetAllAssistantWebhookLogRequest{
		ProjectId: projectId,
		Paginate:  paginate,
		Criterias: criterias,
		Order:     ordering,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantWebhookLog", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantWebhookLog", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}
func (client *assistantServiceClient) GetAssistantWebhookLog(c context.Context,
	auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantWebhookLogRequest) (*assistant_api.GetAssistantWebhookLogResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantWebhookLog(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantWebhookLog", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantWebhookLog %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get GetAssistantWebhookLog %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantWebhookLog", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAssistantWebhook(c context.Context,
	auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantWebhook(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantWebhook", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantWebhook %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get GetAssistantWebhook %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantWebhook", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantWebhook(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantWebhook", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantWebhook %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantWebhook", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) DeleteAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.DeleteAssistantWebhook(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantWebhook", time.Since(start))
		client.logger.Errorf("error while calling DeleteAssistantWebhook %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantWebhook", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) UpdateAssistantWebhook(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantWebhook(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantWebhook", time.Since(start))
		client.logger.Errorf("error while calling UpdateAssistantWebhook %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantWebhook", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAssistantConversation(
	c context.Context,
	auth types.SimplePrinciple,
	assistantRequest *assistant_api.GetAssistantConversationRequest) (*assistant_api.GetAssistantConversationResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantConversation(client.WithAuth(c, auth), assistantRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantConversation", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantConversation %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantConversation", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantAnalysis(ctx context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantAnalysis, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantAnalysis(client.WithAuth(ctx, auth), &assistant_api.GetAllAssistantAnalysisRequest{
		Paginate:    paginate,
		AssistantId: assistantId,
		Criterias:   criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantAnalysis", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantAnalysis", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAssistantAnalysis(c context.Context,
	auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantAnalysis(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantAnalysis", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantAnalysis %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get GetAssistantAnalysis %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantAnalysis", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantAnalysis(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantAnalysis", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantAnalysis %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantAnalysis", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) DeleteAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.DeleteAssistantAnalysis(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantAnalysis", time.Since(start))
		client.logger.Errorf("error while calling DeleteAssistantAnalysis %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantAnalysis", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) UpdateAssistantAnalysis(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantAnalysisRequest) (*assistant_api.GetAssistantAnalysisResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantAnalysis(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantAnalysis", time.Since(start))
		client.logger.Errorf("error while calling UpdateAssistantAnalysis %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantAnalysis", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantTool(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantTool, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantTool(client.WithAuth(c, auth), &assistant_api.GetAllAssistantToolRequest{
		AssistantId: assistantId,
		Paginate:    paginate,
		Criterias:   criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAssistantTool", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAssistantTool", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAssistantTool(c context.Context,
	auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantTool(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantTool", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantTool %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get GetAssistantTool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantTool", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantTool(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantTool", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantTool %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantTool", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) DeleteAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.DeleteAssistantTool(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantTool", time.Since(start))
		client.logger.Errorf("error while calling DeleteAssistantTool %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantTool", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) UpdateAssistantTool(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantTool(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantTool", time.Since(start))
		client.logger.Errorf("error while calling UpdateAssistantTool %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantTool", time.Since(start))
	return res, nil
}

//

func (client *assistantServiceClient) GetAllAssistantKnowledge(c context.Context, auth types.SimplePrinciple, assistantId uint64, criterias []*assistant_api.Criteria, paginate *assistant_api.Paginate) (*assistant_api.Paginated, []*assistant_api.AssistantKnowledge, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantKnowledge(client.WithAuth(c, auth), &assistant_api.GetAllAssistantKnowledgeRequest{
		AssistantId: assistantId,
		Paginate:    paginate,
		Criterias:   criterias,
	})
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAssistantKnowledge", time.Since(start))
		client.logger.Errorf("error while calling to get all assistant %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all assistant %v", err)
	}

	client.logger.Benchmark("Benchmarking: assistantServiceClient.GetAssistantKnowledge", time.Since(start))
	return res.GetPaginated(), res.GetData(), nil
}

func (client *assistantServiceClient) GetAssistantKnowledge(c context.Context,
	auth types.SimplePrinciple, iRequest *assistant_api.GetAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantKnowledge(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantKnowledge", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantKnowledge %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get GetAssistantKnowledge %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantKnowledge", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) CreateAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.CreateAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.CreateAssistantKnowledge(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantKnowledge", time.Since(start))
		client.logger.Errorf("error while calling CreateAssistantKnowledge %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.CreateAssistantKnowledge", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) DeleteAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.DeleteAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.DeleteAssistantKnowledge(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantKnowledge", time.Since(start))
		client.logger.Errorf("error while calling DeleteAssistantKnowledge %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.DeleteAssistantKnowledge", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) UpdateAssistantKnowledge(c context.Context, auth types.SimplePrinciple, iRequest *assistant_api.UpdateAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.UpdateAssistantKnowledge(client.WithAuth(c, auth), iRequest)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantKnowledge", time.Since(start))
		client.logger.Errorf("error while calling UpdateAssistantKnowledge %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.UpdateAssistantKnowledge", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAssistantToolLog(c context.Context, auth types.SimplePrinciple, in *assistant_api.GetAssistantToolLogRequest) (*assistant_api.GetAssistantToolLogResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAssistantToolLog(client.WithAuth(c, auth), in)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantToolLog", time.Since(start))
		client.logger.Errorf("error while calling GetAssistantToolLog %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAssistantToolLog", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantToolLog(ctx context.Context, auth types.SimplePrinciple, in *assistant_api.GetAllAssistantToolLogRequest) (*assistant_api.GetAllAssistantToolLogResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantToolLog(client.WithAuth(ctx, auth), in)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantToolLog", time.Since(start))
		client.logger.Errorf("error while calling GetAllAssistantToolLog %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantToolLog", time.Since(start))
	return res, nil
}

func (client *assistantServiceClient) GetAllAssistantTelemetry(ctx context.Context, auth types.SimplePrinciple, in *assistant_api.GetAllAssistantTelemetryRequest) (*assistant_api.GetAllAssistantTelemetryResponse, error) {
	start := time.Now()
	res, err := client.assistantClient.GetAllAssistantTelemetry(client.WithAuth(ctx, auth), in)
	if err != nil {
		client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantTelemetry", time.Since(start))
		client.logger.Errorf("error while calling GetAllAssistantTelemetry %v", err)
		return nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get tool %v", err)
	}
	client.logger.Benchmark("Benchmarking: assistantClient.GetAllAssistantTelemetry", time.Since(start))
	return res, nil
}
