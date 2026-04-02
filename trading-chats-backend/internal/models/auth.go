package models

import "time"

type Role string

type TenantType string

const (
	RoleAdmin  Role = "admin"
	RoleTenant Role = "tenant"

	TenantTypeSystem TenantType = "system"
	TenantTypeBiz    TenantType = "biz"
)

type Tenant struct {
	ID        string      `bson:"_id,omitempty" json:"id"`
	Name      string      `bson:"name" json:"name"`
	Code      string      `bson:"code" json:"code"`
	Type      TenantType  `bson:"type" json:"type"`
	Status    string      `bson:"status" json:"status"`
	CreatedAt interface{} `bson:"created_at" json:"created_at"`
	UpdatedAt interface{} `bson:"updated_at" json:"updated_at"`
}

type User struct {
	ID           string      `bson:"_id,omitempty" json:"id"`
	TenantID     string      `bson:"tenant_id" json:"tenant_id"`
	Username     string      `bson:"username" json:"username"`
	PasswordHash string      `bson:"password_hash" json:"-"`
	DisplayName  string      `bson:"display_name" json:"display_name"`
	Role         Role        `bson:"role" json:"role"`
	Status       string      `bson:"status" json:"status"`
	LastLoginAt  interface{} `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	CreatedAt    interface{} `bson:"created_at" json:"created_at"`
	UpdatedAt    interface{} `bson:"updated_at" json:"updated_at"`
}

type UserSession struct {
	ID           string      `bson:"_id,omitempty" json:"id"`
	UserID       string      `bson:"user_id" json:"user_id"`
	TenantID     string      `bson:"tenant_id" json:"tenant_id"`
	Role         Role        `bson:"role" json:"role"`
	AccessToken  string      `bson:"access_token" json:"access_token"`
	RefreshToken string      `bson:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time   `bson:"expires_at" json:"expires_at"`
	CreatedAt    interface{} `bson:"created_at" json:"created_at"`
	UpdatedAt    interface{} `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	TenantName  string `json:"tenant_name" binding:"required"`
	TenantCode  string `json:"tenant_code" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthContext struct {
	UserID    string
	TenantID  string
	Role      Role
	Username  string
	SessionID string
}
