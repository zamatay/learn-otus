# Собираем в гошке
FROM golang:1.21.7-alpine as build
LABEL stage=builder

WORKDIR /app

## Кэшируем слои с модулями
COPY integration_test/go.mod .
COPY integration_test/go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o test integration_test/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/test /app/test

CMD ["./test"]