package auth

import "time"

const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"

	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusBanned   = "banned"
)

type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	Bio          *string   `json:"bio"`
	Phone        *string   `json:"phone"`
	AvatarURL    *string   `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Username string `json:"username" binding:"required,min=3,max=80"`
	Email    string `json:"email" binding:"required,email,max=190"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateStaffUserRequest struct {
	Name     string
	Username string
	Email    string
	Password string
	Role     string
}
