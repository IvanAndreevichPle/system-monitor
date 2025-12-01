# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "$LDFLAGS" \
    -o /build/bin/monitor ./cmd/monitor

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "$LDFLAGS" \
    -o /build/bin/client ./cmd/client

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/bin/monitor /app/monitor
COPY --from=builder /build/bin/client /app/client

# Copy example config
COPY --from=builder /build/configs/config.example.yaml /app/configs/config.example.yaml

# Make binaries executable
RUN chmod +x /app/monitor /app/client

# Default command
CMD ["/app/monitor"]

