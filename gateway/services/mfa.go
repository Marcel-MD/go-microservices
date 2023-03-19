package services

import (
	"context"
	"gateway/config"
	"gateway/pb"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MfaService interface {
	GenerateOtp(ctx context.Context, email string) (string, error)
	VerifyOtp(ctx context.Context, email string, otp string) (bool, error)
}

type mfaService struct {
	conn   *grpc.ClientConn
	client pb.MfaServiceClient
}

var (
	mfaOnce sync.Once
	mfaSrv  MfaService
)

func GetMfaService() MfaService {
	mfaOnce.Do(func() {
		cfg := config.GetConfig()

		conn, err := grpc.Dial(cfg.MfaServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to mfa service")
		}

		client := pb.NewMfaServiceClient(conn)

		mfaSrv = &mfaService{
			conn:   conn,
			client: client,
		}
	})

	return mfaSrv
}

func (s *mfaService) Close() error {
	return s.conn.Close()
}

func (s *mfaService) GenerateOtp(ctx context.Context, email string) (string, error) {
	otp, err := s.client.GenerateOtp(ctx, &pb.GenerateOtpRequest{Email: email})
	if err != nil {
		return "", err
	}

	return otp.GetOtp(), nil
}

func (s *mfaService) VerifyOtp(ctx context.Context, email string, otp string) (bool, error) {
	verifyResponse, err := s.client.VerifyOtp(ctx, &pb.VerifyOtpRequest{Email: email, Otp: otp})
	if err != nil {
		return false, err
	}

	return verifyResponse.IsValid, nil
}
