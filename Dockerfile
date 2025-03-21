# Stage 1: Build Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required system packages
RUN apk add --no-cache \
	git \
	upx \
	ca-certificates \
	tzdata

# Download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -extldflags '-static'" \
    -o /app/main .

# Compress binary
RUN upx --best --lzma /app/main

# Stage 2: Minimal runtime image
FROM scratch

# Import certificates, timezone and migrations data
COPY --from=builder --chown=65534:65534 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=65534:65534 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder --chown=65534:65534 /app/migrations /app/migrations

# Copy binary
COPY --from=builder --chown=nobody:nogroup /app/main /main

# Runtime configuration
ENV TZ=UTC
EXPOSE 8080
USER 65534:65534

ENTRYPOINT ["/main"]
