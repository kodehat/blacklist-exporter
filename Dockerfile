FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o blacklist_exporter main.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=builder /app/.env .
COPY --from=builder /app/blacklist_exporter .

EXPOSE 2112

CMD ["/app/blacklist_exporter"]