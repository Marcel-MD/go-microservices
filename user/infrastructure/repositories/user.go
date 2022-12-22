package repositories

import (
	"user/domain"

	"gorm.io/gorm"
)

type IUserRepository interface {
	FindAll() []domain.User
	FindById(id string) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	SearchByEmail(email string) []domain.User
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(user *domain.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindAll() []domain.User {
	var users []domain.User
	r.DB.Find(&users)
	return users
}

func (r *userRepository) FindById(id string) (domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, "id = ?", id).Error

	return user, err
}

func (r *userRepository) SearchByEmail(email string) []domain.User {
	var users []domain.User
	r.DB.Where("email LIKE ?", "%"+email+"%").Find(&users)
	return users
}

func (r *userRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, "email = ?", email).Error

	return user, err
}

func (r *userRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) Update(user *domain.User) error {
	return r.DB.Save(user).Error
}

func (r *userRepository) Delete(user *domain.User) error {
	return r.DB.Delete(user).Error
}
