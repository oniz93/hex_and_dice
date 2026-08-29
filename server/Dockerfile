FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# --- Runtime stage ---
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary and data
COPY --from=builder /server /app/server
COPY --from=builder /app/data /app/data

EXPOSE 8080

ENTRYPOINT ["/app/server"]
