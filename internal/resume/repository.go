package resume

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resume struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	IsDefault      bool      `json:"is_default"`
	AttachmentURLs []string  `json:"attachment_urls"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Resume, error) {
	query := `SELECT id, user_id, title, content::text, is_default, attachment_urls, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM resumes WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list resumes: %w", err)
	}
	defer rows.Close()

	var resumes []Resume
	for rows.Next() {
		var res Resume
		var attachmentURLs []string
		err := rows.Scan(&res.ID, &res.UserID, &res.Title, &res.Content, &res.IsDefault, &attachmentURLs, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan resume: %w", err)
		}
		if attachmentURLs == nil {
			attachmentURLs = []string{}
		}
		res.AttachmentURLs = attachmentURLs
		resumes = append(resumes, res)
	}
	if resumes == nil {
		resumes = []Resume{}
	}
	return resumes, nil
}

func (r *Repository) Create(ctx context.Context, userID uuid.UUID, title, content string) (*Resume, error) {
	// Validate content is valid JSON
	if !json.Valid([]byte(content)) {
		return nil, fmt.Errorf("invalid json content")
	}

	query := `
		INSERT INTO resumes (user_id, title, content)
		VALUES ($1, $2, $3::jsonb)
		RETURNING id, user_id, title, content::text, is_default, attachment_urls, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`
	var res Resume
	var attachmentURLs []string
	err := r.pool.QueryRow(ctx, query, userID, title, content).Scan(
		&res.ID, &res.UserID, &res.Title, &res.Content, &res.IsDefault, &attachmentURLs, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create resume: %w", err)
	}
	if attachmentURLs == nil {
		attachmentURLs = []string{}
	}
	res.AttachmentURLs = attachmentURLs
	return &res, nil
}

func (r *Repository) Update(ctx context.Context, id, userID uuid.UUID, title, content string) error {
	if !json.Valid([]byte(content)) {
		return fmt.Errorf("invalid json content")
	}
	query := `UPDATE resumes SET title = $3, content = $4::jsonb, updated_at = NOW() WHERE id = $1 AND user_id = $2`
	result, err := r.pool.Exec(ctx, query, id, userID, title, content)
	if err != nil {
		return fmt.Errorf("update resume: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM resumes WHERE id = $1 AND user_id = $2`
	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete resume: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (*Resume, error) {
	query := `SELECT id, user_id, title, content::text, is_default, attachment_urls, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM resumes WHERE id = $1 AND user_id = $2`
	var res Resume
	var attachmentURLs []string
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&res.ID, &res.UserID, &res.Title, &res.Content, &res.IsDefault, &attachmentURLs, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("resume not found")
		}
		return nil, fmt.Errorf("get resume: %w", err)
	}
	if attachmentURLs == nil {
		attachmentURLs = []string{}
	}
	res.AttachmentURLs = attachmentURLs
	return &res, nil
}

func (r *Repository) SetDefault(ctx context.Context, id, userID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Unset all defaults for this user
	_, err = tx.Exec(ctx, `UPDATE resumes SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("unset defaults: %w", err)
	}

	// Set the specified resume as default
	result, err := tx.Exec(ctx, `UPDATE resumes SET is_default = true WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("set default: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("resume not found")
	}

	return tx.Commit(ctx)
}

func (r *Repository) AddAttachment(ctx context.Context, id, userID uuid.UUID, url string) error {
	query := `UPDATE resumes SET attachment_urls = array_append(attachment_urls, $3), updated_at = NOW() WHERE id = $1 AND user_id = $2`
	result, err := r.pool.Exec(ctx, query, id, userID, url)
	if err != nil {
		return fmt.Errorf("add attachment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}
