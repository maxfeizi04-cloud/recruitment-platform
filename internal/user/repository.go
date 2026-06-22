package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type HRCertification struct {
	UserID      uuid.UUID `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Position    string    `json:"position"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetProfile(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	query := `SELECT id, phone, role, name, avatar_url FROM users WHERE id = $1`
	var p Profile
	err := r.pool.QueryRow(ctx, query, userID).Scan(&p.ID, &p.Phone, &p.Role, &p.Name, &p.AvatarURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, name, avatarURL string) error {
	query := `UPDATE users SET name = $2, avatar_url = $3, updated_at = NOW() WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, userID, name, avatarURL)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *Repository) SubmitCertification(ctx context.Context, userID uuid.UUID, companyName, position string) error {
	query := `
		INSERT INTO hr_certifications (user_id, company_name, position, status)
		VALUES ($1, $2, $3, 'pending')
		ON CONFLICT (user_id) DO UPDATE SET company_name = $2, position = $3, status = 'pending', updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, query, userID, companyName, position)
	if err != nil {
		return fmt.Errorf("submit certification: %w", err)
	}
	return nil
}

func (r *Repository) GetCertification(ctx context.Context, userID uuid.UUID) (*HRCertification, error) {
	query := `SELECT user_id, company_name, position, status, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM hr_certifications WHERE user_id = $1`
	var c HRCertification
	err := r.pool.QueryRow(ctx, query, userID).Scan(&c.UserID, &c.CompanyName, &c.Position, &c.Status, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // no certification yet
		}
		return nil, fmt.Errorf("get certification: %w", err)
	}
	return &c, nil
}
