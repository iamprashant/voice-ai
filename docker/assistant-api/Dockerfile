# syntax=docker/dockerfile:1
# ---- Build stage ----
FROM rapidaai/rapida-golang:1.25.7-bookworm AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# ---- Copy only what this service needs ----
COPY config/ ./config/
COPY protos/ ./protos/
COPY pkg/ ./pkg/
COPY api/assistant-api ./api/assistant-api
COPY cmd/assistant/ ./cmd/assistant/

# ---- Build Go binary with all CGO flags in one RUN command ----
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    export CGO_CFLAGS="-I/opt/onnxruntime/include -I/usr/local/include -I/opt/azure-speech-sdk/include/c_api" && \
    export CGO_LDFLAGS="-L/opt/onnxruntime/lib -lonnxruntime -L/usr/local/lib -lrnnoise -L/opt/azure-speech-sdk/lib/x64 -lMicrosoft.CognitiveServices.Speech.core" && \
    export LD_LIBRARY_PATH="/opt/onnxruntime/lib:/usr/local/lib:/opt/azure-speech-sdk/lib/x64:$LD_LIBRARY_PATH" && \
    CGO_ENABLED=1 go build -v -o assistant-api ./cmd/assistant/assistant.go

# ---- Runtime stage ----
FROM rapidaai/rapida-debian:bookworm-slim AS runtime

# Install assistant-api specific runtime deps
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && apt-get install -y --no-install-recommends \
    libopus0 libopusfile0

# Copy binary and libraries
COPY --from=builder /app/assistant-api .
COPY --from=builder /opt/onnxruntime /opt/onnxruntime
COPY --from=builder /opt/azure-speech-sdk /opt/azure-speech-sdk
COPY --from=builder /usr/local/lib /usr/local/lib

# Copy migrations and env
COPY api/assistant-api/migrations/ ./api/assistant-api/migrations/
COPY docker/assistant-api/.assistant.env ./env/.assistant.env

USER rapida-app

ENV LD_LIBRARY_PATH="/usr/local/lib:/opt/onnxruntime/lib:/opt/azure-speech-sdk/lib/x64"

EXPOSE 9007

CMD ["./assistant-api"]
