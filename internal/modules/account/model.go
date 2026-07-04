package account

import "github.com/CulipBlue/backend_ednic/internal/modules/auth"

type UpdateProfileRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=120"`
	Username  string  `json:"username" binding:"required,min=3,max=80"`
	Bio       *string `json:"bio" binding:"omitempty,max=500"`
	Phone     *string `json:"phone" binding:"omitempty,max=40"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,max=500"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=72"`
}

type ProfileResponse struct {
	User auth.User `json:"user"`
}
