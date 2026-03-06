package dto

type LoginRequest struct {
	TenantID   *string `json:"tenant_id"`
	Identifier string  `json:"identifier" binding:"required"`
	Password   string  `json:"password" binding:"required"`
}

type GenerateInviteRequest struct {
	TenantID   string `json:"tenant_id" binding:"required"`
	RoleID     string `json:"role_id" binding:"required"`
	Identifier string `json:"identifier" binding:"required"`
}

type AcceptInviteRequest struct {
	InviteToken string `json:"invite_token" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	TenantID   *string `json:"tenant_id"`
	Identifier string  `json:"identifier" binding:"required"`
}

type ResetPasswordRequest struct {
	TenantID    *string `json:"tenant_id"`
	Identifier  string  `json:"identifier" binding:"required"`
	OTP         string  `json:"otp" binding:"required,len=6"`
	NewPassword string  `json:"new_password" binding:"required,min=8"`
}
