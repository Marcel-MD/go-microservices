package services

import (
	"errors"
	"sync"
	"user/auth"
	"user/config"
	"user/repositories"

	"user/domain"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	FindAll() []domain.User
	FindOne(id string) (domain.User, error)
	Register(user domain.User) (domain.User, error)
	Login(email, password string) (string, error)
}

type userService struct {
	repository repositories.IUserRepository
	cfg        config.Config
}

var userOnce sync.Once
var userSrv IUserService

func GetUserService() IUserService {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user service")

		userSrv = &userService{
			repository: repositories.GetUserRepository(),
			cfg:        config.GetConfig(),
		}
	})

	return userSrv
}

func (s *userService) FindAll() []domain.User {
	log.Debug().Msg("Finding all users")

	return s.repository.FindAll()
}

func (s *userService) FindOne(id string) (domain.User, error) {
	log.Debug().Str("id", id).Msg("Finding user")

	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) Register(user domain.User) (domain.User, error) {
	log.Debug().Msg("Registering user")

	existingUser, err := s.repository.FindByEmail(user.Email)
	if err == nil {
		return existingUser, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashedPassword)

	err = s.repository.Create(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (string, error) {
	log.Debug().Msg("Logging in user")

	existingUser, err := s.repository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = s.verifyPassword(password, existingUser.Password)
	if err != nil {
		return "", err
	}

	return auth.Generate(existingUser.Id, s.cfg.TokenLifespan, s.cfg.ApiSecret)
}

func (s *userService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
