package model

import (
	"time"

	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
)

type User struct {
	ID              int64      `bun:"id,pk,autoincrement"`
	Name            string     `bun:"name,notnull"`
	Email           string     `bun:"email,unique,notnull"`
	Password        string     `bun:"password,notnull"`
	EmailVerifiedAt *time.Time `bun:"email_verified_at"`
	CreatedAt       time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Password:        u.Password,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
