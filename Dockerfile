# =============================================================================
# Stage 1: migrate — installs goose for running database migrations
# =============================================================================
FROM golang:1.24-alpine AS migrate

RUN apk add --no-cache netcat-openbsd \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

# =============================================================================
# Stage 2: dev — hot-reload development image using Air
# =============================================================================
FROM golang:1.25-alpine AS dev

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Cache dependency downloads as a separate layer
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# =============================================================================
# Stage 3: builder — compiles the optimised production binary
# =============================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Disable CGO for a fully static binary
ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o /app/cmd-server ./cmd/server

# =============================================================================
# Stage 4: production — minimal runtime image
# =============================================================================
FROM alpine:latest AS production

# Required for HTTPS calls (e.g. Neon TLS)
RUN apk add --no-cache ca-certificates

# Copy compiled binary from builder stage
COPY --from=builder /app/cmd-server /app/cmd-server

# Run as a non-root user for security
RUN adduser -D app \
    && chown app:app /app/cmd-server

USER app

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["/app/cmd-server"]
