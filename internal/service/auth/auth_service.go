package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/akfaiz/go-mailgen"
	"github.com/akfaiz/go-vue-starter-kit/internal/config"
	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/errdefs"
	"github.com/akfaiz/go-vue-starter-kit/internal/validator"
	"github.com/invopop/ctxi18n/i18n"
)

type service struct {
	cfg            config.Config
	userRepo       domain.UserRepository
	userTokenRepo  domain.UserTokenRepository
	passwordHasher domain.PasswordHasher
	jwtManager     domain.JWTManager
	mailer         domain.Mailer
}

func NewService(
	cfg config.Config,
	userRepo domain.UserRepository,
	userTokenRepo domain.UserTokenRepository,
	passwordHasher domain.PasswordHasher,
	jwtManager domain.JWTManager,
	mailer domain.Mailer,
) domain.AuthService {
	return &service{
		cfg:            cfg,
		userRepo:       userRepo,
		userTokenRepo:  userTokenRepo,
		passwordHasher: passwordHasher,
		jwtManager:     jwtManager,
		mailer:         mailer,
	}
}

func (s *service) Register(ctx context.Context, user *domain.User) (*domain.PairToken, error) {
	hashedPassword, err := s.passwordHasher.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, validator.NewError("email", "Email already registered")
		}
		return nil, err
	}

	return s.jwtManager.GeneratePairToken(&domain.JWTClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *service) Login(ctx context.Context, email, password string) (*domain.PairToken, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, validator.NewError("email", i18n.T(ctx, "auth.failed"))
	}

	match, err := s.passwordHasher.Verify(password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, validator.NewError("email", i18n.T(ctx, "auth.failed"))
	}
	claims := &domain.JWTClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return s.jwtManager.GeneratePairToken(claims)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*domain.PairToken, error) {
	claims, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(ctx, claims.ID)
	if err != nil {
		return nil, err
	}
	claims = &domain.JWTClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return s.jwtManager.GeneratePairToken(claims)
}

func (s *service) SendForgotPasswordEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			return validator.NewError("email", i18n.T(ctx, "passwords.user"))
		}
		return err
	}

	token := s.generateRandomString(32)
	hashedToken, err := s.passwordHasher.Hash(token)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(s.cfg.Auth.ResetPasswordExpiration)
	passwordResetToken := &domain.UserToken{
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
		TokenType: domain.TokenTypeResetPassword,
	}
	if err := s.userTokenRepo.Create(ctx, passwordResetToken); err != nil {
		return err
	}

	msg := s.buildEmailResetPassword(user, token)
	if err := s.mailer.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *service) ValidateResetPassword(ctx context.Context, token, email string) error {
	_, err := s.validateResetPassword(ctx, token, email)
	return err
}

func (s *service) ResetPassword(ctx context.Context, token, email, newPassword string) error {
	user, err := s.validateResetPassword(ctx, token, email)
	if err != nil {
		return err
	}

	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user.ID, &domain.UserUpdate{
		Password: omit.From(hashedPassword),
	}); err != nil {
		return err
	}
	_ = s.userTokenRepo.Delete(ctx, user.ID, domain.TokenTypeResetPassword)

	return nil
}

func (s *service) SendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	token := s.generateRandomString(32)
	hashedToken, err := s.passwordHasher.Hash(token)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(s.cfg.Auth.ResetPasswordExpiration)
	passwordResetToken := &domain.UserToken{
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
		TokenType: domain.TokenTypeVerification,
	}
	if err := s.userTokenRepo.Create(ctx, passwordResetToken); err != nil {
		return err
	}
	msg := s.buildEmailVerification(user, token)
	if err := s.mailer.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *service) VerifyEmail(ctx context.Context, token string, userID int64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	prt, err := s.userTokenRepo.FindOne(ctx, user.ID, domain.TokenTypeVerification)
	if err != nil {
		return err
	}

	if time.Now().After(prt.ExpiresAt) {
		return errdefs.ErrBadRequest(i18n.T(ctx, "passwords.token"))
	}

	match, err := s.passwordHasher.Verify(token, prt.Token)
	if err != nil {
		return err
	}
	if !match {
		return errdefs.ErrBadRequest(i18n.T(ctx, "passwords.token"))
	}

	if user.IsVerified() {
		return nil
	}

	if err := s.userRepo.Update(ctx, user.ID, &domain.UserUpdate{
		EmailVerifiedAt: omitnull.From(time.Now()),
	}); err != nil {
		return err
	}
	_ = s.userTokenRepo.Delete(ctx, user.ID, domain.TokenTypeVerification)

	return nil
}

func (s *service) validateResetPassword(ctx context.Context, token, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, errdefs.ErrBadRequest(i18n.T(ctx, "passwords.user"))
		}
		return nil, err
	}

	prt, err := s.userTokenRepo.FindOne(ctx, user.ID, domain.TokenTypeResetPassword)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, errdefs.ErrBadRequest(i18n.T(ctx, "passwords.token"))
		}
		return nil, err
	}

	if time.Now().After(prt.ExpiresAt) {
		return nil, errdefs.ErrBadRequest(i18n.T(ctx, "passwords.token"))
	}

	match, err := s.passwordHasher.Verify(token, prt.Token)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errdefs.ErrBadRequest(i18n.T(ctx, "passwords.token"))
	}

	return user, nil
}

func (s *service) buildEmailResetPassword(user *domain.User, token string) *mailgen.Builder {
	url := s.cfg.App.FrontendBaseURL + "/reset-password?token=" + token + "&email=" + user.Email
	return mailgen.New().
		To(user.Email).
		Subject("Reset Password Notification").
		Name(user.Name).
		Line("You are receiving this email because we received a password reset request for your account.").
		Action("Reset Password", url).
		Linef("This password reset link will expire in %d minutes.", 60).
		Line("If you did not request a password reset, no further action is required.")
}

func (s *service) buildEmailVerification(user *domain.User, token string) *mailgen.Builder {
	url := s.cfg.App.FrontendBaseURL + "/verify-email?token=" + token + "&email=" + user.Email
	return mailgen.New().
		To(user.Email).
		Subject("Verify Email Address").
		Name(user.Name).
		Line("Please click the button below to verify your email address.").
		Action("Verify Email Address", url).
		Line("If you did not create an account, no further action is required.")
}

func (s *service) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var max = big.NewInt(int64(len(charset)))

	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, max)
		if err != nil {
			return ""
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}
