package rpc

import (
	"context"
	"net"
	"sync"
	"user/config"
	"user/models"
	"user/pb"
	"user/services"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			userService: services.GetUserService(),
			mailService: services.GetMailService(),
		}

		l, err := net.Listen("tcp", cfg.Port)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen.")
		}

		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, server)

		listener = l
		srv = s
	})

	return srv, listener
}

type server struct {
	pb.UserServiceServer
	userService services.UserService
	mailService services.MailService
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.UserId, error) {

	user := models.User{
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
	}

	createdUser, err := s.userService.Register(user)
	if err != nil {
		return &pb.UserId{}, err
	}

	s.mailService.SendWelcomeMail(createdUser)

	return &pb.UserId{Id: createdUser.Id}, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Token, error) {

	token, err := s.userService.Login(in.Email, in.Password)
	if err != nil {
		return &pb.Token{}, err
	}

	return &pb.Token{Token: token}, nil
}

func (s *server) GetAll(ctx context.Context, in *empty.Empty) (*pb.Users, error) {

	users := s.userService.FindAll()

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        user.Id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.Users{Users: pbUsers}, nil
}

func (s *server) Get(ctx context.Context, in *pb.UserId) (*pb.User, error) {

	user, err := s.userService.FindOne(in.Id)
	if err != nil {
		return &pb.User{}, err
	}

	return &pb.User{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
