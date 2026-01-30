// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package integration_api

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

// StreamChatBidirectional handles bidirectional streaming chat with persistent connection.
//
// This method:
// 1. Authenticates the client once at the beginning.
// 2. Keeps the connection open, receiving multiple ChatRequest messages.
// 3. For each message, processes it through the LLM caller.
// 4. Sends responses back through the same stream.
// 5. Continues until the client closes the stream or an error occurs.
//
// Advantages:
// - Single persistent connection for multiple messages
// - No reconnection overhead
// - Real-time bidirectional communication
// - Efficient for conversational AI
//
// Parameters:
// - context: The context for the request, used for authentication and cancellation.
// - providerName: A string identifier for the provider, used for logging.
// - callerFactory: A function to create a new LLM caller for each request with its credential.
// - stream: The bidirectional gRPC stream for receiving requests and sending responses.
//
// Returns an error if authentication fails or if there's an issue during streaming.
func (iApi *integrationApi) StreamChatBidirectional(
	context context.Context,
	providerName string,
	callerFactory func(*protos.Credential) internal_callers.LargeLanguageCaller,
	stream grpc.BidiStreamingServer[protos.ChatRequest, protos.ChatResponse],
) error {
	// Authenticate once at the beginning
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(context)
	if !isAuthenticated || !iAuth.HasProject() {
		iApi.logger.Errorf("unauthenticated request for bidirectional stream chat")
		return status.Error(codes.Unauthenticated, "Please provide valid service credentials to perform invoke.")
	}

	iApi.logger.Infof("Bidirectional stream chat opened for provider: %s", providerName)

	// Keep connection open and process multiple requests
	for {
		// Receive next chat request from client
		irRequest, err := stream.Recv()
		if err == io.EOF {
			// Client closed the stream gracefully
			iApi.logger.Infof("Client closed bidirectional stream for provider: %s", providerName)
			return nil
		}
		if err != nil {
			iApi.logger.Errorf("Error receiving from bidirectional stream: %v", err)
			return status.Errorf(codes.Internal, "Error receiving chat request from stream: %v", err)
		}

		if irRequest == nil {
			iApi.logger.Warnf("Received nil request from bidirectional stream")
			continue
		}

		// Generate unique request ID for this message
		uuID := iApi.RequestId()
		if irRequest.AdditionalData == nil {
			irRequest.AdditionalData = map[string]string{}
		}

		// Populate request metadata
		irRequest.AdditionalData["provider_name"] = providerName
		model, ok := irRequest.ModelParameters["model.name"]
		if ok {
			mdl, err := utils.AnyToString(model)
			if err == nil {
				irRequest.AdditionalData["model_name"] = mdl
			}
		}

		modelID, ok := irRequest.ModelParameters["model.id"]
		if ok {
			mdlID, err := utils.AnyToString(modelID)
			if err == nil {
				irRequest.AdditionalData["model_id"] = mdlID
			}
		}

		source, ok := utils.GetClientSource(context)
		if ok {
			irRequest.AdditionalData["source"] = source.Get()
		}

		clientEnv, ok := utils.GetClientEnvironment(context)
		if ok {
			irRequest.AdditionalData["env"] = clientEnv.Get()
		}

		clientRegion, ok := utils.GetClientRegion(context)
		if ok {
			irRequest.AdditionalData["region"] = clientRegion.Get()
		}

		// Create a new LLM caller for this request with its credential
		llmCaller := callerFactory(irRequest.GetCredential())

		// Process the chat completion request
		err = llmCaller.StreamChatCompletion(
			stream.Context(),
			irRequest.GetConversations(),
			internal_callers.NewChatOptions(
				uuID,
				irRequest,
				iApi.PreHook(stream.Context(), iAuth, irRequest, uuID, providerName),
				iApi.PostHook(stream.Context(), iAuth, irRequest, uuID, providerName),
			),
			func(rID string, content *protos.Message) error {
				return stream.Send(&protos.ChatResponse{
					Success:   true,
					RequestId: rID,
					Data:      content,
				})
			},
			func(rID string, content *protos.Message, mtx []*protos.Metric) error {
				return stream.Send(&protos.ChatResponse{
					Success:   true,
					RequestId: rID,
					Metrics:   mtx,
					Data:      content,
				})
			},
			func(rID string, err error) {
				stream.Send(&protos.ChatResponse{
					Success:   false,
					Code:      400,
					RequestId: rID,
					Error: &protos.Error{
						ErrorCode:    uint64(400),
						ErrorMessage: err.Error(),
						HumanMessage: err.Error(),
					},
				})
			},
		)

		// If there's an error during processing, send it and continue (don't close stream)
		if err != nil {
			iApi.logger.Warnf("Error processing chat request in bidirectional stream: %v", err)
			stream.Send(&protos.ChatResponse{
				Success:   false,
				Code:      500,
				RequestId: irRequest.GetRequestId(),
				Error: &protos.Error{
					ErrorCode:    500,
					ErrorMessage: err.Error(),
					HumanMessage: "Internal server error processing your request",
				},
			})
			// Continue to next request instead of closing stream
		}
	}
}

func (iApi *integrationApi) Chat(
	c context.Context,
	irRequest *protos.ChatRequest,
	tag string,
	caller internal_callers.LargeLanguageCaller,
) (*protos.ChatResponse, error) {
	iApi.logger.Infof("Chat from grpc with provider %s", tag)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated || !iAuth.HasProject() {
		iApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[protos.ChatResponse](
			errors.New("unauthenticated request for chat"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	uuID := iApi.RequestId()
	if irRequest.AdditionalData == nil {
		irRequest.AdditionalData = map[string]string{}
	}

	irRequest.AdditionalData["provider_name"] = tag
	model, ok := irRequest.ModelParameters["model.name"]
	if ok {
		mdl, err := utils.AnyToString(model)
		if err == nil {
			irRequest.AdditionalData["model_name"] = mdl
		}
	}

	modelID, ok := irRequest.ModelParameters["model.id"]
	if ok {
		mdlID, err := utils.AnyToString(modelID)
		if err == nil {
			irRequest.AdditionalData["model_id"] = mdlID
		}
	}
	source, ok := utils.GetClientSource(c)
	if ok {
		irRequest.AdditionalData["source"] = source.Get()
	}

	clientEnv, ok := utils.GetClientEnvironment(c)
	if ok {
		irRequest.AdditionalData["env"] = clientEnv.Get()
	}

	clientRegion, ok := utils.GetClientRegion(c)
	if ok {
		irRequest.AdditionalData["region"] = clientRegion.Get()
	}

	completions, metrics, err := caller.GetChatCompletion(
		c,
		irRequest.GetConversations(),
		internal_callers.NewChatOptions(
			uuID,
			irRequest,
			iApi.PreHook(c, iAuth, irRequest, uuID, tag),
			iApi.PostHook(c, iAuth, irRequest, uuID, tag),
		),
	)
	if err != nil {
		return utils.Error[protos.ChatResponse](err, err.Error())
	}
	return &protos.ChatResponse{
		Code:    200,
		Success: true,
		Data:    completions,
		Metrics: metrics,
	}, nil
}
