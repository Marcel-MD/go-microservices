package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"user/config"
	"user/models"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type IMailService interface {
	Send(ctx context.Context, mail models.Mail) error
	SendAsync(mail models.Mail)
	Close() error

	SendWelcomeMail(user models.User)
}

type mailService struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

var mailOnce sync.Once
var mailSrv IMailService

const queueName = "mail"

func GetMailService() IMailService {
	mailOnce.Do(func() {
		log.Info().Msg("Initializing mail service")

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

		mailSrv = &mailService{
			conn:  conn,
			ch:    ch,
			queue: q,
		}
	})

	return mailSrv
}

func (s *mailService) Close() error {
	err := s.ch.Close()
	if err != nil {
		return err
	}

	return s.conn.Close()
}

func (s *mailService) Send(ctx context.Context, mail models.Mail) error {
	body, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	err = s.ch.PublishWithContext(
		ctx,
		"",           // exchange
		s.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func (s *mailService) SendAsync(mail models.Mail) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Send(ctx, mail)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
	}()
}

func (s *mailService) SendWelcomeMail(user models.User) {
	mail := models.Mail{
		To:      []string{user.Email},
		Subject: "Welcome to the service",
		Body:    fmt.Sprintf("Hello <b>%s</b>! Welcome to the service", user.FirstName),
	}

	s.SendAsync(mail)
}
