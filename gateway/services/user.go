package services

import (
	"context"
	"gateway/config"
	"gateway/models"
	"gateway/pb"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	Close() error
	Get(ctx context.Context, id string) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Register(ctx context.Context, user models.RegisterUser) (string, error)
	Login(ctx context.Context, user models.LoginUser) (string, error)
}

type userService struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

var (
	userOnce sync.Once
	userSrv  UserService
)

func GetUserService() UserService {
	userOnce.Do(func() {
		cfg := config.GetConfig()

		conn, err := grpc.Dial(cfg.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to user service")
		}

		client := pb.NewUserServiceClient(conn)

		userSrv = &userService{
			conn:   conn,
			client: client,
		}
	})

	return userSrv
}

func (s *userService) Close() error {
	return s.conn.Close()
}

func (s *userService) Get(ctx context.Context, id string) (models.User, error) {
	user, err := s.client.Get(ctx, &pb.UserId{Id: id})
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}, nil
}

func (s *userService) List(ctx context.Context) ([]models.User, error) {
	users, err := s.client.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	var result []models.User
	for _, user := range users.Users {
		result = append(result, models.User{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.AsTime(),
			UpdatedAt: user.UpdatedAt.AsTime(),
		})
	}

	return result, nil
}

func (s *userService) Register(ctx context.Context, user models.RegisterUser) (string, error) {
	resp, err := s.client.Register(ctx, &pb.RegisterRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

func (s *userService) Login(ctx context.Context, user models.LoginUser) (string, error) {
	resp, err := s.client.Login(ctx, &pb.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}
