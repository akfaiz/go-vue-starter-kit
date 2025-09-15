package handler

import (
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/handler/dto"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/middleware/auth"
	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/errdefs"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	pairToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}

	res := dto.NewResponse(200, dto.NewTokenResponse(pairToken), "Login successful")
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	user := req.ToDomain()
	ctx := c.Request().Context()
	if err := h.authService.Register(ctx, user); err != nil {
		return err
	}

	res := dto.NewMessage(201, "User registered successfully")
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	pairToken, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return err
	}

	res := dto.NewResponse(200, dto.NewTokenResponse(pairToken))
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) SendForgotPasswordEmail(c echo.Context) error {
	var req dto.SendForgotPasswordEmailRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.authService.SendForgotPasswordEmail(ctx, req.Email); err != nil {
		return err
	}

	res := dto.NewMessage(200, i18n.T(ctx, "passwords.sent"))
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) ValidateResetPassword(c echo.Context) error {
	var req dto.ValidateResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.authService.ValidateResetPassword(ctx, req.Token, req.Email); err != nil {
		return err
	}

	res := dto.NewMessage(200, "Reset password token is valid")
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.authService.ResetPassword(ctx, req.Token, req.Email, req.NewPassword); err != nil {
		return err
	}

	res := dto.NewMessage(200, "Password has been reset successfully")
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) SendVerificationEmail(c echo.Context) error {
	user := auth.GetUser(c)
	if user == nil {
		return errdefs.ErrUnauthorized()
	}
	ctx := c.Request().Context()
	if err := h.authService.SendVerificationEmail(ctx, user.Email); err != nil {
		return err
	}

	res := dto.NewMessage(200, "If the email is registered, a verification link has been sent")
	return c.JSON(res.Status, res)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var req dto.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.authService.VerifyEmail(ctx, req.Token, req.Email); err != nil {
		return err
	}

	res := dto.NewMessage(200, "Email has been verified successfully")
	return c.JSON(res.Status, res)
}
