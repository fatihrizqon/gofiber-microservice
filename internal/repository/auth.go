package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Login(username string) (entity.User, error)
}

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepository {
	return &AuthRepository{Db: Db}
}

func (r *AuthRepository) Login(email string) (entity.User, error) {
	var entity entity.User
	if err := r.Db.Where("email = ?", email).First(&entity).Error; err != nil {
		return entity, errors.New("credentials does not matches our record")
	}
	return entity, nil
}
