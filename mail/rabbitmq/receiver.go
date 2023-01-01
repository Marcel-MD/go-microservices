package rabbitmq

import (
	"context"
	"encoding/json"
	"mail/config"
	"mail/models"
	"mail/services"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type IReceiver interface {
	Listen(ctx context.Context)
	Close() error
}

type receiver struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
	msgs  <-chan amqp.Delivery

	mailService services.IMailService
}

var receiverOnce sync.Once
var receive IReceiver

const queueName = "mail"

func GetReceiver() IReceiver {
	receiverOnce.Do(func() {
		cfg := config.GetConfig()

		conn, err := amqp.Dial(cfg.RabbitMQUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open a channel")
		}

		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to declare a queue")
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to register a consumer")
		}

		receive = &receiver{
			conn:  conn,
			ch:    ch,
			queue: q,
			msgs:  msgs,

			mailService: services.GetMailService(),
		}
	})

	return receive
}

func (r *receiver) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}

	return r.conn.Close()
}

func (r *receiver) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-r.msgs:
			go r.processMessage(d)
		}
	}
}

func (r *receiver) processMessage(msg amqp.Delivery) {
	log.Info().Msg("Processing message")

	var mail models.Mail

	err := json.Unmarshal(msg.Body, &mail)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message")
		return
	}

	r.mailService.Send(mail)
}
