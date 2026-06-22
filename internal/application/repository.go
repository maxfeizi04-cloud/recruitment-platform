package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	ID        uuid.UUID `json:"id"`
	JobID     uuid.UUID `json:"job_id"`
	UserID    uuid.UUID `json:"user_id"`
	ResumeID  uuid.UUID `json:"resume_id"`
	Status    string    `json:"status"`
	JobTitle  string    `json:"job_title,omitempty"`
	UserName  string    `json:"user_name,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, userID, jobID, resumeID uuid.UUID) (*Application, error) {
	query := `
		INSERT INTO job_applications (job_id, user_id, resume_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (job_id, user_id) DO UPDATE SET resume_id = $3, status = 'pending', updated_at = NOW()
		RETURNING id, job_id, user_id, resume_id, status,
			to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`
	var a Application
	err := r.pool.QueryRow(ctx, query, jobID, userID, resumeID).Scan(
		&a.ID, &a.JobID, &a.UserID, &a.ResumeID, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}
	return &a, nil
}

func (r *Repository) ListByCandidate(ctx context.Context, userID uuid.UUID) ([]Application, error) {
	query := `
		SELECT a.id, a.job_id, a.user_id, a.resume_id, a.status, j.title,
			to_char(a.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(a.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM job_applications a JOIN jobs j ON a.job_id = j.id
		WHERE a.user_id = $1 ORDER BY a.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query applications: %w", err)
	}
	defer rows.Close()

	var apps []Application
	for rows.Next() {
		var a Application
		if err := rows.Scan(&a.ID, &a.JobID, &a.UserID, &a.ResumeID, &a.Status, &a.JobTitle, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		apps = append(apps, a)
	}
	if apps == nil {
		apps = []Application{}
	}
	return apps, nil
}

func (r *Repository) ListByHR(ctx context.Context, hrUserID uuid.UUID) ([]Application, error) {
	query := `
		SELECT a.id, a.job_id, a.user_id, a.resume_id, a.status, j.title, u.name,
			to_char(a.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(a.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM job_applications a JOIN jobs j ON a.job_id = j.id JOIN users u ON a.user_id = u.id
		WHERE j.hr_user_id = $1 ORDER BY a.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, hrUserID)
	if err != nil {
		return nil, fmt.Errorf("query applications: %w", err)
	}
	defer rows.Close()

	var apps []Application
	for rows.Next() {
		var a Application
		if err := rows.Scan(&a.ID, &a.JobID, &a.UserID, &a.ResumeID, &a.Status, &a.JobTitle, &a.UserName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		apps = append(apps, a)
	}
	if apps == nil {
		apps = []Application{}
	}
	return apps, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id, hrUserID uuid.UUID, status string) error {
	query := `
		UPDATE job_applications a SET status = $3, updated_at = NOW()
		FROM jobs j WHERE a.job_id = j.id AND a.id = $1 AND j.hr_user_id = $2
	`
	result, err := r.pool.Exec(ctx, query, id, hrUserID, status)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("application not found or not owned by you")
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	query := `
		SELECT a.id, a.job_id, a.user_id, a.resume_id, a.status, j.title, u.name,
			to_char(a.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(a.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM job_applications a JOIN jobs j ON a.job_id = j.id JOIN users u ON a.user_id = u.id
		WHERE a.id = $1
	`
	var a Application
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.JobID, &a.UserID, &a.ResumeID, &a.Status, &a.JobTitle, &a.UserName, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("application not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}
	return &a, nil
}
