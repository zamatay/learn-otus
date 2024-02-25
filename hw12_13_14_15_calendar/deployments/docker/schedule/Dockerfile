# Собираем в гошке
FROM golang:1.21.7-alpine as build
LABEL stage=builder

ENV BIN_FILE schedule

WORKDIR /app

## Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ${BIN_FILE} cmd/schedule/main.go

# На выходе тонкий образ
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/${BIN_FILE} /app/${BIN_FILE}
COPY ./configs/config.yaml /app/config.yaml


ENTRYPOINT ["./schedule", "-config", "/app/config.yaml"]
