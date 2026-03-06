package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tenant struct {
	Base
	Name       string `gorm:"type:varchar(255);not null"`
	Domain     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	AuthConfig string `gorm:"type:jsonb"`
	IsActive   bool   `gorm:"default:true"`
	Users      []User `gorm:"foreignKey:TenantID"`
	Roles      []Role `gorm:"foreignKey:TenantID"`
}

type IdentifierType string

const (
	IdentifierEmail IdentifierType = "EMAIL"
	IdentifierPhone IdentifierType = "PHONE"
)

type UserStatus string

const (
	StatusInvited   UserStatus = "INVITED"
	StatusActive    UserStatus = "ACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
)

type User struct {
	Base
	TenantID          *uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_tenant_identifier"`
	Tenant            *Tenant
	IdentifierType    IdentifierType `gorm:"type:varchar(20);not null"`
	PrimaryIdentifier string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_identifier"`
	PasswordHash      string         `gorm:"type:varchar(255);not null"`
	FirstName         string         `gorm:"type:varchar(100);not null"`
	LastName          string         `gorm:"type:varchar(100)"`
	Status            UserStatus     `gorm:"type:varchar(20);default:'INVITED'"`

	Roles []Role `gorm:"many2many:user_roles;"`
}

type Permission struct {
	ID          string    `gorm:"type:varchar(100);primaryKey"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Role struct {
	Base
	Name        string       `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_role"`
	Description string       `gorm:"type:text"`
	IsSystem    bool         `gorm:"default:false"`
	TenantID    *uuid.UUID   `gorm:"type:uuid;uniqueIndex:idx_tenant_role"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type InviteStatus string

const (
	InvitePending  InviteStatus = "PENDING"
	InviteAccepted InviteStatus = "ACCEPTED"
	InviteExpired  InviteStatus = "EXPIRED"
)

type Invitation struct {
	Base
	TenantID   uuid.UUID    `gorm:"type:uuid;not null;index"`
	InviterID  uuid.UUID    `gorm:"type:uuid;not null"`
	Identifier string       `gorm:"type:varchar(255);not null"`
	RoleID     uuid.UUID    `gorm:"type:uuid;not null"`
	TokenHash  string       `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status     InviteStatus `gorm:"type:varchar(20);default:'PENDING'"`
	ExpiresAt  time.Time    `gorm:"not null"`
}

type UserSession struct {
	Base
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	OpaqueTokenHash string    `gorm:"type:varchar(255);not null;uniqueIndex"`

	IPAddress  string `gorm:"type:varchar(45)"`
	UserAgent  string `gorm:"type:text"`
	DeviceType string `gorm:"type:varchar(50)"`
	Browser    string `gorm:"type:varchar(50)"`
	OS         string `gorm:"type:varchar(50)"`
	Location   string `gorm:"type:varchar(255)"`

	IsActive     bool      `gorm:"default:true;index"`
	ExpiresAt    time.Time `gorm:"not null"`
	LastActiveAt time.Time `gorm:"not null"`
}
