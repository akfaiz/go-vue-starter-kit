package dto

import (
	"time"

	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
)

type ProfileResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required" label:"Name"`
	Email string `json:"email" validate:"required|email" label:"Email"`
}

type ChangePasswordRequest struct {
	CurrentPassword         string `json:"current_password" validate:"required" label:"Current Password"`
	NewPassword             string `json:"new_password" validate:"required|min:8" label:"New Password"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required|eq_field:NewPassword" label:"Confirm New Password"`
}

type DeleteAccountRequest struct {
	Password string `json:"password" validate:"required" label:"Password"`
}

func NewProfileResponse(user *domain.User) *ProfileResponse {
	return &ProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
