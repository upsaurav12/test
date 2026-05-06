# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

# Install certificates for HTTPS calls inside the container.
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Cache dependency downloads as a separate layer.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Inject build-time version information.
ARG VERSION=dev
ARG BUILD_TIME

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
      -X main.Version=${VERSION} \
      -X 'main.BuildTime=${BUILD_TIME}'" \
    -o /app/bin/server ./cmd/main.go

# ── Final stage ───────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/server /server

# Expose the application port.
EXPOSE 8080

ENTRYPOINT ["/server"]
