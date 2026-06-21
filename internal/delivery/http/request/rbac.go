package request

type AssignRoleRequest struct {
	RoleId string `validate:"required,uuid" json:"role_id"`
}

type AssignPermissionRequest struct {
	PermissionId string `validate:"required,uuid" json:"permission_id"`
}
