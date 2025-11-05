package workflow_routers

import (
	"github.com/gin-gonic/gin"
	assistantApi "github.com/rapidaai/api/assistant-api/assistant"
	assistantDeploymentApi "github.com/rapidaai/api/assistant-api/assistant-deployment"
	assistantTalkApi "github.com/rapidaai/api/assistant-api/talk"
	"github.com/rapidaai/config"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	workflow_api "github.com/rapidaai/protos"
	"google.golang.org/grpc"
)

func AssistantApiRoute(
	Cfg *config.AppConfig,
	S *grpc.Server,
	Logger commons.Logger,
	Postgres connectors.PostgresConnector,
	Redis connectors.RedisConnector,
	Opensearch connectors.OpenSearchConnector,
) {
	workflow_api.RegisterAssistantServiceServer(S,
		assistantApi.NewAssistantGRPCApi(Cfg,
			Logger,
			Postgres,
			Redis,
			Opensearch,
			Opensearch,
		))
}

func AssistantDeploymentApiRoute(Cfg *config.AppConfig,
	S *grpc.Server,
	Logger commons.Logger,
	Postgres connectors.PostgresConnector) {
	workflow_api.RegisterAssistantDeploymentServiceServer(S,
		assistantDeploymentApi.NewAssistantDeploymentGRPCApi(Cfg,
			Logger,
			Postgres,
		))
}

func AssistantConversationApiRoute(
	Cfg *config.AppConfig,
	S *grpc.Server,
	Logger commons.Logger,
	Postgres connectors.PostgresConnector,
	Redis connectors.RedisConnector,
	Opensearch connectors.OpenSearchConnector,
) {
	workflow_api.RegisterTalkServiceServer(S,
		assistantTalkApi.NewConversationGRPCApi(Cfg,
			Logger,
			Postgres,
			Redis,
			Opensearch,
			Opensearch,
		))
}

func TalkCallbackApiRoute(
	cfg *config.AppConfig, engine *gin.Engine, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector) {
	apiv1 := engine.Group("v1/talk")
	talkRpcApi := assistantTalkApi.NewConversationApi(cfg,
		logger,
		postgres,
		redis,
		opensearch,
		opensearch,
	)
	{
		// for incomming call
		// https://integral-presently-cub.ngrok-free.app/v1/talk/exotel/call/2200665081979600896?x-api-key=cc0d7b49cd51480d46c1d68dd37fed5ecd6663a5f5285a0ebc7ffb7a9cd3e0b2
		apiv1.GET("/exotel/call/:assistantId", talkRpcApi.ExotelCallReciever)
		apiv1.GET("/exotel/usr/:assistantId/:identifier/:conversationId/:authorization/:x-auth-id/:x-project-id", talkRpcApi.ExotelCallTalker)
		apiv1.GET("/exotel/prj/:assistantId/:identifier/:conversationId/:x-api-key", talkRpcApi.ExotelCallTalker)

		// 2193965165215481856/+14582073109/2206320653408141312/bca387c4e5cb8fd4cbcaeb194389216959c14dcaee4069d76c52093f5f571171
		// /v1/talk/twilio/stream/2200665081979600896/+14582073109/2206374850140831744

		// /v1/talk/twilio/usr/2219166493587800064/+6596522466/2224218665107062784/61c814ba2a3868574e53860537bb4bc03a9bd1305a822800d5f0ee0c1206ac5c/2021822161534058496
		// only for debugger
		apiv1.GET("/twilio/call/:assistantId", talkRpcApi.PhoneCallReciever)
		apiv1.GET("/twilio/usr/:assistantId/:identifier/:conversationId/:authorization/:x-auth-id/:x-project-id", talkRpcApi.TwilioCallTalker)
		apiv1.GET("/twilio/prj/:assistantId/:identifier/:conversationId/:x-api-key", talkRpcApi.TwilioCallTalker)
		apiv1.POST("/twilio/whatsapp/:assistantToken", talkRpcApi.WhatsappReciever)

		//
		apiv1.GET("/vonage/call/:assistantId", talkRpcApi.PhoneCallReciever)
		apiv1.GET("/vonage/usr/:assistantId/:identifier/:conversationId/:authorization/:x-auth-id/:x-project-id", talkRpcApi.VonageCallTalker)
		apiv1.GET("/vonage/prj/:assistantId/:identifier/:conversationId/:x-api-key", talkRpcApi.VonageCallTalker)

	}
}

func ConversationApiRoute(
	cfg *config.AppConfig, engine *gin.Engine, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector) {
	apiv1 := engine.Group("v1/conversation")
	talkRpcApi := assistantTalkApi.NewConversationApi(cfg,
		logger,
		postgres,
		redis,
		opensearch,
		opensearch,
	)
	apiv1.POST("/create-phone-call", talkRpcApi.ExotelCallReciever)
	apiv1.POST("/create-bulk-phone-call", talkRpcApi.ExotelCallReciever)

}
