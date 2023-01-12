package services

import (
	"context"
	"encoding/json"
	"fmt"
	"mfa/config"
	"mfa/models"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
)

type IMailService interface {
	Send(ctx context.Context, mail models.Mail) error
	SendAsync(mail models.Mail)
	Close() error

	SendOtpMail(email, otp string)
}

type mailService struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
}

var (
	mailOnce sync.Once
	mailSrv  IMailService
)

const queueName = "mail"

func GetMailService() IMailService {
	mailOnce.Do(func() {
		log.Info().Msg("Initializing mail service")

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

		publisher, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsLogging,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create publisher")
		}

		mailSrv = &mailService{
			conn:      conn,
			publisher: publisher,
		}
	})

	return mailSrv
}

func (s *mailService) Close() error {
	s.publisher.Close()
	return s.conn.Close()
}

func (s *mailService) Send(ctx context.Context, mail models.Mail) error {
	body, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	return s.publisher.Publish(
		body,
		[]string{queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
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

func (s *mailService) SendOtpMail(email, otp string) {
	mail := models.Mail{
		To:      []string{email},
		Subject: "Verification Code",
		Body:    fmt.Sprintf("Your verification code is <strong>%s</strong>.", otp),
	}

	s.SendAsync(mail)
}
