package services

import (
	"context"
	"math/rand"
	"mfa/config"
	"mfa/repositories"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
)

type IOtpService interface {
	Generate(ctx context.Context, email string) (string, error)
	Verify(ctx context.Context, email, otp string) (bool, error)
}

type otpService struct {
	cfg        config.Config
	repository repositories.IOtpRepository
}

var (
	otpOnce sync.Once
	otpSrv  IOtpService
)

func GetOtpService() IOtpService {
	otpOnce.Do(func() {
		log.Info().Msg("Initializing otp service")

		otpSrv = &otpService{
			cfg:        config.GetConfig(),
			repository: repositories.GetOtpRepository(),
		}
	})

	return otpSrv
}

func (s *otpService) Generate(ctx context.Context, email string) (string, error) {
	num := 100000 + rand.Intn(800000)
	otp := strconv.Itoa(num)

	err := s.repository.Set(ctx, email, otp, s.cfg.OtpExpiry)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (s *otpService) Verify(ctx context.Context, email, otp string) (bool, error) {
	actualOtp, err := s.repository.Get(ctx, email)
	if err != nil {
		return false, err
	}

	return actualOtp == otp, nil
}
