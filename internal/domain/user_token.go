package domain

import (
	"context"
	"time"
)

type UserTokenRepository interface {
	Create(ctx context.Context, token *UserToken) error
	FindOne(ctx context.Context, userID int64, tokenType TokenType) (*UserToken, error)
	Delete(ctx context.Context, userID int64, tokenType TokenType) error
}

type UserToken struct {
	ID        int64
	UserID    int64
	Token     string
	TokenType TokenType
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TokenType string

const (
	TokenTypeVerification  TokenType = "verification"
	TokenTypeResetPassword TokenType = "reset_password"
)
