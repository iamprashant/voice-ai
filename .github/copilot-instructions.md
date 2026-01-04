# Rapida Voice AI Platform - AI Coding Guidelines

## Architecture Overview

Rapida is a Go-based microservices platform for voice AI orchestration. Key components:

- **Backend Services**: web-api (port 9001), assistant-api (9007), integration-api (9004), endpoint-api (9005), document-api (9010)
- **Shared Packages**: `pkg/` contains reusable components (clients, models, commons, connectors)
- **Frontend**: React/TypeScript UI (port 3000) with Tailwind CSS
- **Infrastructure**: Docker Compose with PostgreSQL, Redis, OpenSearch
- **Communication**: REST APIs via Gin, internal gRPC with protobuf definitions

## Development Workflow

- **Setup**: Run `make up-all` to start all services via Docker Compose
- **Build**: Use `make build-all` for container builds; individual services with `make build-web`
- **Logs**: Monitor with `make logs-web` or `make logs-all`
- **Testing**: Execute tests with `go test ./...` or `yarn test` in UI
- **Dependencies**: Go modules managed with `go mod`; Node packages with Yarn

## Key Patterns & Conventions

### Go Backend Patterns

- **Logging**: Use `commons.Logger` interface (zap-based) - inject via constructor
- **Configuration**: Structs with `mapstructure` tags, validated with `validator/v10`
- **REST Clients**: Implement `APIClient` interface; use `APIResponse.ToJSON()`, `.ToMap()`, `.ToString()` for responses
- **Error Handling**: Return errors with context; use `fmt.Errorf("failed to X: %w", err)`
- **Routing**: Gin groups under `/v1` for REST; gRPC services registered with protobuf
- **Database**: GORM models with migrations in `migrations/` directories

### Frontend Patterns

- **Styling**: Tailwind CSS with custom build via `yarn build:css`
- **Scripts**: Use `yarn start:dev` for concurrent Tailwind watch + Craco dev server
- **Linting**: ESLint with `yarn lint:fix`; Prettier for formatting

### Code Examples

- **REST Client Usage**:
  ```go
  client := rest.NewRestClient(logger, cfg, "https://api.example.com")
  resp, err := client.Get(ctx, "/endpoint", params, headers)
  if err != nil { return err }
  data, err := resp.ToMap()
  ```
- **Logger Injection**:
  ```go
  type Service struct {
      logger commons.Logger
  }
  func NewService(logger commons.Logger) *Service {
      return &Service{logger: logger}
  }
  ```
- **Config Struct**:
  ```go
  type Config struct {
      APIKey string `mapstructure:"api_key" validate:"required"`
  }
  ```

## Integration Points

- **OAuth**: Multiple providers (Google, GitHub, etc.) configured in `OAuth2Config`
- **AI Services**: Anthropic, OpenAI, Cohere clients for LLM interactions
- **Voice Processing**: STT/TTS via Google Cloud Speech, Azure Cognitive Services
- **Storage**: PostgreSQL for relational data, Redis for caching, OpenSearch for search

## File Organization

- `api/*/`: Service-specific code (handlers, routers)
- `pkg/`: Shared utilities and models
- `protos/`: Generated gRPC code from `.proto` files
- `ui/src/`: React components and logic
- `docker/*/`: Service-specific Dockerfiles and configs

Focus on microservice boundaries; use gRPC for inter-service communication, REST for external APIs.</content>
<parameter name="filePath">.github/copilot-instructions.md
