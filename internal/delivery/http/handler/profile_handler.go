package handler

import (
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/handler/dto"
	"github.com/akfaiz/go-vue-starter-kit/internal/delivery/http/middleware/auth"
	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/errdefs"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	userService domain.UserService
}

func NewProfileHandler(userService domain.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()

	claims := auth.GetUser(c)
	if claims == nil {
		return errdefs.ErrUnauthorized()
	}
	user, err := h.userService.FindByID(ctx, claims.ID)
	if err != nil {
		return err
	}

	res := dto.NewResponse(200, dto.NewProfileResponse(user))
	return c.JSON(res.Status, res)
}

func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	claims := auth.GetUser(c)
	if claims == nil {
		return errdefs.ErrUnauthorized()
	}
	user, err := h.userService.FindByID(ctx, claims.ID)
	if err != nil {
		return err
	}
	user.Name = req.Name
	user.Email = req.Email
	if err := h.userService.UpdateProfile(ctx, claims.ID, user); err != nil {
		return err
	}
	res := dto.NewResponse(200, dto.NewProfileResponse(user))
	return c.JSON(res.Status, res)
}

func (h *ProfileHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	claims := auth.GetUser(c)
	if claims == nil {
		return errdefs.ErrUnauthorized()
	}
	if err := h.userService.ChangePassword(ctx, claims.ID, req.CurrentPassword, req.NewPassword); err != nil {
		return err
	}
	res := dto.NewMessage(200, "Password changed successfully")
	return c.JSON(res.Status, res)
}

func (h *ProfileHandler) DeleteAccount(c echo.Context) error {
	var req dto.DeleteAccountRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	claims := auth.GetUser(c)
	if claims == nil {
		return errdefs.ErrUnauthorized()
	}
	if err := h.userService.Delete(ctx, claims.ID, req.Password); err != nil {
		return err
	}
	res := dto.NewMessage(200, "Account deleted successfully")
	return c.JSON(res.Status, res)
}
