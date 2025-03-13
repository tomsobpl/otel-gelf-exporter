# Stage 1: Builder
FROM golang:1.23-bookworm AS builder

ARG OTEL_VERSION=0.119.0

WORKDIR /build
RUN go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}

COPY .docker/otel_collector/config-builder.yaml .
COPY . ./gelfexporter
RUN CGO_ENABLED=0 builder --config=config-builder.yaml

# Stage 2: Final Image
FROM cgr.dev/chainguard/static:latest

WORKDIR /app

COPY --from=builder /build/otelcol-custom .
COPY .docker/otel_collector/config.yaml .

EXPOSE 4317/tcp 4318/tcp 13133/tcp
CMD ["/app/otelcol-custom", "--config=config.yaml"]
