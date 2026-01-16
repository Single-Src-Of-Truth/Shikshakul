package domain

import "time"

const (
	UserStatusInvited   = "INVITED"
	UserStatusActive    = "ACTIVE"
	UserStatusSuspended = "SUSPENDED"
)

type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Status       string    `json:"status"`
	IsMFAEnabled bool      `json:"is_mfa_enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InviteRequest struct {
	Email    string `json:"email" validate:"required,email"`
	RoleName string `json:"role_name" validate:"required"`
}

type AcceptInviteRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
