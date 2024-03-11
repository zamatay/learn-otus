# Собираем в гошке
FROM golang:1.21.7-alpine as build
LABEL stage=builder

WORKDIR /app

## Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o calendar cmd/main.go

# На выходе тонкий образ
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/calendar /app/calendar
COPY ./configs/config.yaml /app/config.yaml

EXPOSE 44044 8080

ENTRYPOINT ["./calendar", "-config=/app/config.yaml"]
