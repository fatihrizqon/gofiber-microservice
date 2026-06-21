package service

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IRbacService interface {
	AssignRole(userId string, req request.AssignRoleRequest) error
	RevokeRole(userId string, roleId string) error
	AssignPermission(roleId string, req request.AssignPermissionRequest) error
	RevokePermission(roleId string, permissionId string) error
}

type RbacService struct {
	RbacRepo repository.IRbacRepository
	UserRepo repository.IUserRepository
	Validate *validator.Validate
}

func NewRbacService(rbacRepo repository.IRbacRepository, userRepo repository.IUserRepository, validate *validator.Validate) IRbacService {
	return &RbacService{
		RbacRepo: rbacRepo,
		UserRepo: userRepo,
		Validate: validate,
	}
}

func (s *RbacService) AssignRole(userId string, req request.AssignRoleRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	_, err = s.UserRepo.FindById(parsedUserId)
	if err != nil {
		return errors.New("user not found")
	}

	exists, err := s.RbacRepo.RoleExists(req.RoleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	return s.RbacRepo.AssignRoleToUser(userId, req.RoleId)
}

func (s *RbacService) RevokeRole(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	_, err = s.UserRepo.FindById(parsedUserId)
	if err != nil {
		return errors.New("user not found")
	}

	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	return s.RbacRepo.RevokeRoleFromUser(userId, roleId)
}

func (s *RbacService) AssignPermission(roleId string, req request.AssignPermissionRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	existsPerm, err := s.RbacRepo.PermissionExists(req.PermissionId)
	if err != nil || !existsPerm {
		return errors.New("permission not found")
	}

	return s.RbacRepo.AssignPermissionToRole(roleId, req.PermissionId)
}

func (s *RbacService) RevokePermission(roleId string, permissionId string) error {
	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	existsPerm, err := s.RbacRepo.PermissionExists(permissionId)
	if err != nil || !existsPerm {
		return errors.New("permission not found")
	}

	return s.RbacRepo.RevokePermissionFromRole(roleId, permissionId)
}
