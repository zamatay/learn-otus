package rabbit

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
)

type funcClose = func()

type Rabbited interface {
	SendMessage(ctx context.Context, text string)
}

type Rabbit struct {
	connection *amqp.Connection
	chanel     *amqp.Channel
	queue      amqp.Queue
	closers    []funcClose
}

func NewRabbit(login string, password string, url string) (*Rabbit, error) {
	r := &Rabbit{}
	_, err := r.Connect(login, password, url)
	if err != nil {
		return nil, err
	}

	r.CreateChannel()
	r.CreateQueue("work")
	return r, nil
}

func (r *Rabbit) Connect(login string, password string, url string) (*amqp.Connection, error) {
	var err error
	r.closers = make([]func(), 3)
	r.connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", login, password, url)) // Создаем подключение к RabbitMQ
	if err != nil {
		logger.Logger().Error("не смог подключиться к RabbitMQ", "Error", err.Error())
		return nil, nil
	}

	fnClose := func() {
		_ = r.connection.Close()
	}

	r.closers = append(r.closers, fnClose)

	return r.connection, nil
}

func (r *Rabbit) CreateChannel() {
	var err error
	r.chanel, err = r.connection.Channel()
	if err != nil {
		logger.Logger().Fatal(fmt.Sprintf("не смог создать канал RabbitMQ\n%s", err.Error()))
	}

	fnClose := func() {
		_ = r.chanel.Close() // Закрываем канал в случае удачной попытки открытия
	}

	r.closers = append(r.closers, fnClose)
}

func (r *Rabbit) CreateQueue(nameQueue string) {
	var err error
	r.queue, err = r.chanel.QueueDeclare(
		nameQueue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		logger.Logger().Fatal(fmt.Sprintf("failed to declare a queue. Error: %s", err.Error()))
	}
}

func (r *Rabbit) SendMessage(ctx context.Context, text string) {
	if err := r.chanel.PublishWithContext(ctx,
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(text),
		}); err != nil {
		logger.Logger().Error(fmt.Sprintf("failed to declare a queue. Error: %s", err.Error()))
	}
}

func (r Rabbit) Close() error {
	for _, fn := range r.closers {
		fn()
	}
	return nil
}
