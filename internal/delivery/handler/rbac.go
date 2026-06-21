package handler

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
)

type RbacHandler struct {
	IRbacService service.IRbacService
}

func NewRbacHandler(serv service.IRbacService) *RbacHandler {
	return &RbacHandler{IRbacService: serv}
}

func (h *RbacHandler) AssignRole(ctx fiber.Ctx) error {
	userId := ctx.Params("id")
	req := request.AssignRoleRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IRbacService.AssignRole(userId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully assigned to user.",
	})
}

func (h *RbacHandler) RevokeRole(ctx fiber.Ctx) error {
	userId := ctx.Params("id")
	roleId := ctx.Params("role_id")

	if err := h.IRbacService.RevokeRole(userId, roleId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully revoked from user.",
	})
}

func (h *RbacHandler) AssignPermission(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	req := request.AssignPermissionRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IRbacService.AssignPermission(roleId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully assigned to role.",
	})
}

func (h *RbacHandler) RevokePermission(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	permissionId := ctx.Params("permission_id")

	if err := h.IRbacService.RevokePermission(roleId, permissionId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully revoked from role.",
	})
}
