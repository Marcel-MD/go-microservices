package kafka

import (
	"context"
	"encoding/json"
	"mail/config"
	"mail/models"
	"mail/services"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type IReader interface {
	ReadMessages(ctx context.Context)
}

type reader struct {
	reader      *kafka.Reader
	mailService services.IMailService
}

var readerOnce sync.Once
var r IReader

func GetReader() IReader {
	readerOnce.Do(func() {
		cfg := config.GetConfig()

		r = &reader{
			reader: kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{cfg.KafkaUrl},
				Topic:     cfg.KafkaTopic,
				Partition: cfg.KafkaPartition,
				MinBytes:  10e3, // 10KB
				MaxBytes:  10e6, // 10MB
			}),
			mailService: services.GetMailService(),
		}
	})

	return r
}

func (r *reader) ReadMessages(ctx context.Context) {
	for {
		m, err := r.reader.ReadMessage(ctx)
		if err != nil {
			break
		}

		go r.processMessage(m)
	}

	if err := r.reader.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close reader")
	}
}

func (r *reader) processMessage(m kafka.Message) {
	log.Debug().Msgf("Processing message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))

	var mail models.Mail

	if err := json.Unmarshal(m.Value, &mail); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message")
		return
	}

	r.mailService.Send(mail)
}
