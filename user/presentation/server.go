package presentation

import (
	"context"
	"user/application"
	"user/domain"
	"user/pb"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UserServiceServer
	service application.IUserService
}

func NewServer(service application.IUserService) pb.UserServiceServer {
	return &server{service: service}
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
