package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	query := `SELECT id, phone, role, name, avatar_url FROM users WHERE phone = $1`

	var user User
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Phone, &user.Role, &user.Name, &user.AvatarURL,
	)
	if err == pgx.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("query user by phone: %w", err)
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, phone, role string) (*User, error) {
	query := `
		INSERT INTO users (phone, role)
		VALUES ($1, $2)
		RETURNING id, phone, role, name, avatar_url
	`

	var user User
	err := r.pool.QueryRow(ctx, query, phone, role).Scan(
		&user.ID, &user.Phone, &user.Role, &user.Name, &user.AvatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &user, nil
}

func (r *Repository) FindOrCreate(ctx context.Context, phone, role string) (*User, error) {
	user, err := r.FindByPhone(ctx, phone)
	if err == nil {
		return user, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	return r.Create(ctx, phone, role)
}
