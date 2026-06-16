package repository

import (
	"gorm.io/gorm"
)

type IRbacRepository interface {
	GetUserPermissions(userId string) ([]string, error)
}

type RbacRepository struct {
	Db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) IRbacRepository {
	return &RbacRepository{Db: db}
}

func (r *RbacRepository) GetUserPermissions(userId string) ([]string, error) {
	var permissions []string

	// Ambil semua name permission berdasarkan relasi many-to-many user -> role -> permission
	err := r.Db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userId).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
