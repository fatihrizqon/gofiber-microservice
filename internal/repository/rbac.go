package repository

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRbacRepository interface {
	GetUserPermissions(userId string) ([]string, error)
	AssignRoleToUser(userId string, roleId string) error
	RevokeRoleFromUser(userId string, roleId string) error
	AssignPermissionToRole(roleId string, permissionId string) error
	RevokePermissionFromRole(roleId string, permissionId string) error
	RoleExists(roleId string) (bool, error)
	PermissionExists(permissionId string) (bool, error)
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

func (r *RbacRepository) AssignRoleToUser(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	user := entity.User{Id: parsedUserId}
	role := entity.Role{Id: parsedRoleId}
	return r.Db.Model(&user).Association("Roles").Append(&role)
}

func (r *RbacRepository) RevokeRoleFromUser(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	user := entity.User{Id: parsedUserId}
	role := entity.Role{Id: parsedRoleId}
	return r.Db.Model(&user).Association("Roles").Delete(&role)
}

func (r *RbacRepository) AssignPermissionToRole(roleId string, permissionId string) error {
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	parsedPermissionId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	role := entity.Role{Id: parsedRoleId}
	permission := entity.Permission{Id: parsedPermissionId}
	return r.Db.Model(&role).Association("Permissions").Append(&permission)
}

func (r *RbacRepository) RevokePermissionFromRole(roleId string, permissionId string) error {
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	parsedPermissionId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	role := entity.Role{Id: parsedRoleId}
	permission := entity.Permission{Id: parsedPermissionId}
	return r.Db.Model(&role).Association("Permissions").Delete(&permission)
}

func (r *RbacRepository) RoleExists(roleId string) (bool, error) {
	parsedId, err := uuid.Parse(roleId)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.Db.Model(&entity.Role{}).Where("id = ?", parsedId).Count(&count).Error
	return count > 0, err
}

func (r *RbacRepository) PermissionExists(permissionId string) (bool, error) {
	parsedId, err := uuid.Parse(permissionId)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.Db.Model(&entity.Permission{}).Where("id = ?", parsedId).Count(&count).Error
	return count > 0, err
}


