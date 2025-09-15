package user

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/validator"
)

type service struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
}

func NewService(
	userRepo domain.UserRepository,
	passwordHasher domain.PasswordHasher,
) domain.UserService {
	return &service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (s *service) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *service) UpdateProfile(ctx context.Context, id int64, user *domain.User) error {
	oldUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	update := &domain.UserUpdate{
		Name:  omit.From(user.Name),
		Email: omit.From(user.Email),
	}
	if oldUser.Email != user.Email {
		var t *time.Time
		update.EmailVerifiedAt = omitnull.FromPtr(t)
	}

	return s.userRepo.Update(ctx, id, update)
}

func (s *service) ChangePassword(ctx context.Context, id int64, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	match, err := s.passwordHasher.Verify(currentPassword, user.Password)
	if err != nil {
		return err
	}
	if !match {
		return validator.NewError("current_password", "Current password is incorrect")
	}
	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}
	update := &domain.UserUpdate{
		Password: omit.From(hashedPassword),
	}
	return s.userRepo.Update(ctx, id, update)
}

func (s *service) Delete(ctx context.Context, id int64, password string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	match, err := s.passwordHasher.Verify(password, user.Password)
	if err != nil {
		return err
	}
	if !match {
		return validator.NewError("password", "Password is incorrect")
	}
	return s.userRepo.Delete(ctx, id)
}
