package rabbitmq

import (
	"encoding/json"
	"mail/config"
	"mail/models"
	"mail/services"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
)

type IConsumer interface {
	Close() error
}

type consumer struct {
	conn        *rabbitmq.Conn
	consumer    *rabbitmq.Consumer
	mailService services.IMailService
}

var consumerOnce sync.Once
var consume IConsumer

const queueName = "mail"

func GetConsumer() IConsumer {
	consumerOnce.Do(func() {
		log.Info().Msg("Initializing consumer")

		cfg := config.GetConfig()

		const retries = 5
		var conn *rabbitmq.Conn
		var err error

		for i := 0; i < retries; i++ {
			conn, err = rabbitmq.NewConn(
				cfg.RabbitMQUrl,
				rabbitmq.WithConnectionOptionsLogging,
			)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to connect to RabbitMQ, retrying in 3 seconds...")
				time.Sleep(3 * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
		}

		c := &consumer{
			conn:        conn,
			mailService: services.GetMailService(),
		}

		cons, err := rabbitmq.NewConsumer(
			conn,
			c.handleDelivery,
			queueName,
			rabbitmq.WithConsumerOptionsQueueDurable,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create consumer")
		}

		c.consumer = cons

		consume = c
	})

	return consume
}

func (c *consumer) Close() error {
	c.consumer.Close()
	return c.conn.Close()
}

func (c *consumer) handleDelivery(d rabbitmq.Delivery) rabbitmq.Action {
	log.Info().Msg("Processing delivery")

	var mail models.Mail

	err := json.Unmarshal(d.Body, &mail)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message")
		return rabbitmq.Ack
	}

	c.mailService.Send(mail)

	return rabbitmq.Ack
}
