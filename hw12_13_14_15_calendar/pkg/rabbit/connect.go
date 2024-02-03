package rabbit

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"time"
)

type funcClose = func()

type ProduceRabbited interface {
	SendMessage(ctx context.Context, text string)
}

type ConsumerRabbited interface {
	GetMessage(ctx context.Context) error
}

type Rabbit struct {
	ReadInterval time.Duration
	connection   *amqp.Connection
	chanel       *amqp.Channel
	queue        amqp.Queue
	closers      []funcClose
	queueName    string
}

func (r *Rabbit) GetMessage(ctx context.Context) error {
	messages, err := r.chanel.Consume(r.queueName, "", true, false, false, false, nil)
	if err != nil {
		logger.Logger().Error("Ошибка при получении сообщения", "Error", err.Error())
		return err
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			logger.Logger().Info("received a message: %s", message.Body)
		}
	}()

	logger.Logger().Info(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
	return nil
}

func NewRabbit(config *configs.Broker) (*Rabbit, error) {
	r := &Rabbit{}
	r.queueName = "work"
	_, err := r.Connect(config.Login, config.Password, config.Url)
	if err != nil {
		return nil, err
	}
	r.ReadInterval = config.ReadInterval
	r.CreateChannel()
	r.CreateQueue(r.queueName)
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
