package services

import (
	"errors"
	"testing"
	"time"
	"user/config"
	"user/models"

	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	srv := getUserService()

	user := models.User{
		Email:     "test@mail.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
	}

	user, err := srv.Register(user)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	}

	if user.Password == "password" {
		t.Errorf("Password not hashed")
	}

	_, err = srv.Register(user)
	if err == nil {
		t.Errorf("No error registering existing user")
	}
}

func TestLogin(t *testing.T) {
	srv := getUserService()

	user := models.User{
		Email:     "test@mail.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
	}

	user, err := srv.Register(user)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	}

	token, err := srv.Login(user.Email, "password")
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	}

	if token == "" {
		t.Errorf("Token is empty")
	}

	_, err = srv.Login("wrong@mail.com", "password")
	if err == nil {
		t.Errorf("No error logging in with wrong email")
	}

	_, err = srv.Login(user.Email, "wrongpassword")
	if err == nil {
		t.Errorf("No error logging in with wrong password")
	}
}

func getUserService() IUserService {
	return &userService{
		repository: &userRepositoryMock{
			store: make(map[string]models.User),
		},
		cfg: config.Config{
			ApiSecret:     "SecretSecretSecret",
			TokenLifespan: 24 * time.Hour,
		},
	}
}

type userRepositoryMock struct {
	store map[string]models.User
}

func (r *userRepositoryMock) FindAll() []models.User {
	users := make([]models.User, 0, len(r.store))
	for _, user := range r.store {
		users = append(users, user)
	}
	return users
}

func (r *userRepositoryMock) FindById(id string) (models.User, error) {
	user, ok := r.store[id]
	if !ok {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepositoryMock) FindByEmail(email string) (models.User, error) {
	user, ok := r.store[email]
	if !ok {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepositoryMock) Create(user *models.User) error {
	user.Id = uuid.New().String()
	r.store[user.Id] = *user
	r.store[user.Email] = *user
	return nil
}

func (r *userRepositoryMock) Update(user *models.User) error {
	r.store[user.Id] = *user
	r.store[user.Email] = *user
	return nil
}

func (r *userRepositoryMock) Delete(user *models.User) error {
	delete(r.store, user.Id)
	delete(r.store, user.Email)
	return nil
}
