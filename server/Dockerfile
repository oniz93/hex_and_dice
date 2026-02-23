FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# --- Runtime stage ---
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary and data
COPY --from=builder /server /app/server
COPY --from=builder /app/data /app/data

EXPOSE 8080

ENTRYPOINT ["/app/server"]
