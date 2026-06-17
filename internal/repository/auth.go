package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Register(user entity.User) (entity.User, error)
	Login(username string) (entity.User, error)
}

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepository {
	return &AuthRepository{Db: Db}
}

func (r *AuthRepository) Register(user entity.User) (entity.User, error) {
	var existingUser entity.User
	if !errors.Is(r.Db.Where("username = ?", user.Username).First(&existingUser).Error, gorm.ErrRecordNotFound) {
		return entity.User{}, errors.New("username already taken")
	}
	if !errors.Is(r.Db.Where("email = ?", user.Email).First(&existingUser).Error, gorm.ErrRecordNotFound) {
		return entity.User{}, errors.New("email already taken")
	}
	if err := r.Db.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) Login(email string) (entity.User, error) {
	var entity entity.User
	if err := r.Db.Where("email = ?", email).First(&entity).Error; err != nil {
		return entity, errors.New("credentials does not matches our record")
	}
	return entity, nil
}
