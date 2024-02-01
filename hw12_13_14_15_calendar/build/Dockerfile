# Собираем в гошке
FROM golang:1.2 as build

ENV BIN_FILE /opt/calendar/calendar-app
ENV BIN_FILE_SCHEDULER /opt/calendar/calendar-scheduler
ENV BIN_FILE_SENDER /opt/calendar/calendar-sender
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN gen mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 gen build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/main.go
RUN CGO_ENABLED=0 gen build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_SCHEDULER} cmd/main.go
RUN CGO_ENABLED=0 gen build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_SENDER} cmd/main.go

# На выходе тонкий образ
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE /opt/calendar/calendar-app
ENV BIN_FILE_SCHEDULER /opt/calendar/calendar-scheduler
ENV BIN_FILE_SENDER /opt/calendar/calendar-sender
COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY --from=build ${BIN_FILE_SCHEDULER} ${BIN_FILE_SCHEDULER}
COPY --from=build ${BIN_FILE_SENDER} ${BIN_FILE_SENDER}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
CMD ${BIN_FILE_SCHEDULER} -config ${CONFIG_FILE}
CMD ${BIN_FILE_SENDER} -config ${CONFIG_FILE}