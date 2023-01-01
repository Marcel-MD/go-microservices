package services

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type IMailService interface {
}

type mailService struct {
}

var mailOnce sync.Once
var mailSrv IMailService

func GetMailService() IMailService {
	mailOnce.Do(func() {
		log.Info().Msg("Initializing mail service")

		mailSrv = &mailService{}
	})

	return mailSrv
}
