# arm64
# FROM golang:1.23 AS builder 

# amd64
FROM --platform=linux/amd64 golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

# Сбор бинарника для Linux amd64
RUN GOARCH=amd64 GOOS=linux go build -o /bin/ticker_pulse_bot ./cmd/main.go

# arm64
# FROM alpine:latest

# amd64
FROM --platform=linux/amd64 alpine:latest

RUN apk --no-cache add ca-certificates libc6-compat

COPY --from=builder /bin/ticker_pulse_bot /bin/ticker_pulse_bot

COPY .env /bin/.env

RUN chmod +x /bin/ticker_pulse_bot

ENV ENV_FILE=/bin/.env

RUN ls -l /bin/ticker_pulse_bot

CMD ["/bin/sh", "-c", "/bin/ticker_pulse_bot"]
