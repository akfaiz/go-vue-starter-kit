package dto

import "github.com/akfaiz/go-vue-starter-kit/internal/domain"

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required" label:"Name"`
	Email                string `json:"email" validate:"required|email" label:"Email"`
	Password             string `json:"password" validate:"required|min:6" label:"Password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required|eq_field:Password" label:"Confirm Password"`
}

func (r *RegisterRequest) ToDomain() *domain.User {
	return &domain.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required|email" label:"Email"`
	Password string `json:"password" validate:"required" label:"Password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" label:"Refresh Token"`
}

type ValidateResetPasswordRequest struct {
	Token string `json:"token" validate:"required" label:"Token"`
	Email string `json:"email" validate:"required|email" label:"Email"`
}

type SendForgotPasswordEmailRequest struct {
	Email string `json:"email" validate:"required|email" label:"Email"`
}

type ResetPasswordRequest struct {
	Token                string `json:"token" validate:"required" label:"Token"`
	Email                string `json:"email" validate:"required|email" label:"Email"`
	Password             string `json:"password" validate:"required|min:8" label:"Password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required|eq_field:Password" label:"Confirm Password"`
}

type VerifyEmailRequest struct {
	Token  string `json:"token" validate:"required" label:"Token"`
	UserID int64  `json:"user_id" validate:"required" label:"User ID"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenResponse(token *domain.PairToken) *TokenResponse {
	return &TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}
