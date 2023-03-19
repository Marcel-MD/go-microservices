package rpc

import (
	"context"
	"mfa/config"
	"mfa/pb"
	"mfa/services"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	once     sync.Once
	srv      *grpc.Server
	listener net.Listener
)

func GetServer() (*grpc.Server, net.Listener) {
	once.Do(func() {
		log.Info().Msg("Initializing gRPC server")

		cfg := config.GetConfig()

		server := &server{
			userService: services.GetOtpService(),
			mailService: services.GetMailService(),
		}

		l, err := net.Listen("tcp", cfg.Port)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen.")
		}

		s := grpc.NewServer()
		pb.RegisterMfaServiceServer(s, server)

		listener = l
		srv = s
	})

	return srv, listener
}

type server struct {
	pb.MfaServiceServer
	userService services.OtpService
	mailService services.MailService
}

func (s *server) GenerateOtp(ctx context.Context, in *pb.GenerateOtpRequest) (*pb.OtpResponse, error) {

	otp, err := s.userService.Generate(ctx, in.Email)
	if err != nil {
		return &pb.OtpResponse{}, err
	}

	s.mailService.SendOtpMail(in.Email, otp)

	return &pb.OtpResponse{Otp: otp}, nil
}

func (s *server) VerifyOtp(ctx context.Context, in *pb.VerifyOtpRequest) (*pb.VerifyResponse, error) {

	isValid, err := s.userService.Verify(ctx, in.Email, in.Otp)
	if err != nil {
		return &pb.VerifyResponse{}, err
	}

	return &pb.VerifyResponse{IsValid: isValid}, nil
}
