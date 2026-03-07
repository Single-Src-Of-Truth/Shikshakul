package dto

type CreateTenantRequest struct {
	Name       string `json:"name" binding:"required"`
	Domain     string `json:"domain" binding:"required"`
	AuthConfig string `json:"auth_config"`
}

type TenantResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	IsActive bool   `json:"is_active"`
}

type UpdateTenantRequest struct {
	Name       *string `json:"name"`
	Domain     *string `json:"domain"`
	AuthConfig *string `json:"auth_config"`
	IsActive   *bool   `json:"is_active"`
}

type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1"`
}

type RoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsSystem    bool     `json:"is_system"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateRoleRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

type UserProfileResponse struct {
	ID             string   `json:"id"`
	TenantID       *string  `json:"tenant_id,omitempty"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Identifier     string   `json:"identifier"`
	IdentifierType string   `json:"identifier_type"`
	Status         string   `json:"status"`
	Roles          []string `json:"roles"`
	Permissions    []string `json:"permissions"`
}

type SessionResponse struct {
	IPAddress    string `json:"ip_address"`
	Browser      string `json:"browser"`
	OS           string `json:"os"`
	DeviceType   string `json:"device_type"`
	Location     string `json:"location"`
	LastActiveAt string `json:"last_active_at"`
	IsCurrent    bool   `json:"is_current"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
}

type PermissionDTO struct {
	ID          string `json:"id" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type BulkCreatePermissionsRequest struct {
	Permissions []PermissionDTO `json:"permissions" binding:"required,min=1"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE SUSPENDED"`
}

type ChangeUserRoleRequest struct {
	RoleID string `json:"role_id" binding:"required,uuid"`
}
