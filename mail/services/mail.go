package services

import (
	"fmt"
	"mail/config"
	"mail/models"
	"net/smtp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type MailService interface {
	Send(mail models.Mail)
}

type mailService struct {
	senderName string
	from       string
	addr       string
	auth       smtp.Auth
}

var (
	mailOnce sync.Once
	mailSrv  MailService
)

func GetMailService() MailService {
	mailOnce.Do(func() {
		log.Info().Msg("Initializing mail service")

		cfg := config.GetConfig()

		addr := cfg.SmtpHost + ":" + cfg.SmtpPort
		auth := smtp.PlainAuth("", cfg.SmtpEmail, cfg.SmtpPassword, cfg.SmtpHost)

		log.Info().Str("senderName", cfg.SenderName).Msg("Mail service initialized")

		mailSrv = &mailService{
			from:       cfg.SmtpEmail,
			addr:       addr,
			auth:       auth,
			senderName: cfg.SenderName,
		}
	})
	return mailSrv
}

func (s *mailService) Send(mail models.Mail) {
	log.Debug().Msg("Sending mail")
	msg := s.buildMail(mail)

	err := smtp.SendMail(s.addr, s.auth, s.from, mail.To, msg)
	if err != nil {
		log.Error().Err(err).Msg("Error sending mail")
	}
}

func (s *mailService) buildMail(mail models.Mail) []byte {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", s.senderName)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return []byte(msg)
}
