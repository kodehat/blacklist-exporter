FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build

FROM alpine:3.21.3

LABEL org.opencontainers.image.authors='dev@codehat.de' \
      org.opencontainers.image.url='https://github.com/kodehat/blacklist-exporter' \
      org.opencontainers.image.documentation='https://github.com/kodehat/blacklist-exporter' \
      org.opencontainers.image.source='https://github.com/kodehat/blacklist-exporter' \
      org.opencontainers.image.vendor='kodehat'

WORKDIR /opt

# "curl" is added only for Docker healthchecks!
RUN apk --no-cache add ca-certificates tzdata curl && \
    update-ca-certificates && \
    adduser -D -H nonroot

COPY --from=builder --chown=nonroot:nonroot --chmod=550 /app/blacklist-exporter ./blacklist-exporter

EXPOSE 2112/tcp

USER nonroot:nonroot

ENTRYPOINT [ "/opt/prometheus-storagebox-exporter" ]