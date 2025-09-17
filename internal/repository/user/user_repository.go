package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/akfaiz/go-vue-starter-kit/internal/domain"
	"github.com/akfaiz/go-vue-starter-kit/internal/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) domain.UserRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	m := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	_, err := r.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		var pgError pgdriver.Error
		if errors.As(err, &pgError) {
			if pgError.IntegrityViolation() && strings.Contains(pgError.Error(), "uk_users_email") {
				return domain.ErrEmailAlreadyExists
			}
		}
		return err
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *repository) Update(ctx context.Context, id int64, data *domain.UserUpdate) error {
	query := r.db.NewUpdate().Model((*model.User)(nil)).Where("id = ?", id)
	if data.Name.IsValue() {
		query = query.Set("name = ?", data.Name)
	}
	if data.Email.IsValue() {
		query = query.Set("email = ?", data.Email)
	}
	if data.Password.IsValue() {
		query = query.Set("password = ?", data.Password)
	}
	if data.EmailVerifiedAt.IsValue() {
		if data.EmailVerifiedAt.IsNull() {
			query = query.Set("email_verified_at = NULL")
		} else {
			t := data.EmailVerifiedAt.MustGet()
			query = query.Set("email_verified_at = ?", t)
		}
	}
	query = query.Set("updated_at = NOW()")
	_, err := query.Exec(ctx)
	if err != nil {
		var pgError pgdriver.Error
		if errors.As(err, &pgError) {
			if pgError.IntegrityViolation() && strings.Contains(pgError.Error(), "uk_users_email") {
				return domain.ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model(&model.User{ID: id}).WherePK().Exec(ctx)
	return err
}
