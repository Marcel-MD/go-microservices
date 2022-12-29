package rpc

import (
	"context"
	"net"
	"user/config"
	"user/domain"
	"user/pb"
	"user/services"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Run() {
	log.Info().Msg("Initializing gRPC server")

	cfg := config.GetConfig()

	server := &server{
		service: services.GetUserService(),
	}

	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen.")
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server)

	log.Info().Msg("Starting server")

	if err := s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve.")
	}
}

type server struct {
	pb.UserServiceServer
	service services.IUserService
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.UserId, error) {

	user := domain.User{
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
	}

	createdUser, err := s.service.Register(user)
	if err != nil {
		return &pb.UserId{}, err
	}

	return &pb.UserId{Id: createdUser.Id}, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Token, error) {

	token, err := s.service.Login(in.Email, in.Password)
	if err != nil {
		return &pb.Token{}, err
	}

	return &pb.Token{Token: token}, nil
}

func (s *server) List(ctx context.Context, in *empty.Empty) (*pb.UserList, error) {

	users := s.service.FindAll()

	var userList []*pb.User
	for _, user := range users {
		userList = append(userList, &pb.User{
			Id:        user.Id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.UserList{Users: userList}, nil
}

func (s *server) Get(ctx context.Context, in *pb.UserId) (*pb.User, error) {

	user, err := s.service.FindOne(in.Id)
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
