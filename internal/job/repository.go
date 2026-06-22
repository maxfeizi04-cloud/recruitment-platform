package job

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Job struct {
	ID           uuid.UUID `json:"id"`
	HRUserID     uuid.UUID `json:"hr_user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements"` // JSON string
	SalaryRange  string    `json:"salary_range"` // JSON string
	Location     string    `json:"location"`     // JSON string
	Status       string    `json:"status"`
	ExpireAt     string    `json:"expire_at"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type SearchParams struct {
	Query  string
	City   string
	Salary string
	Limit  int
	Offset int
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func scanJob(row pgx.Row) (*Job, error) {
	var j Job
	err := row.Scan(&j.ID, &j.HRUserID, &j.Title, &j.Description, &j.Requirements, &j.SalaryRange, &j.Location, &j.Status, &j.ExpireAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Job, int, error) {
	countQuery := `SELECT COUNT(*) FROM jobs WHERE status = 'active'`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count jobs: %w", err)
	}

	query := `SELECT id, hr_user_id, title, description, requirements::text, salary_range::text, location::text, status, to_char(expire_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM jobs WHERE status = 'active' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, *j)
	}
	if jobs == nil {
		jobs = []Job{}
	}
	return jobs, total, nil
}

func (r *Repository) Search(ctx context.Context, params SearchParams) ([]Job, int, error) {
	conditions := []string{"status = 'active'"}
	args := []interface{}{}
	argIdx := 1

	if params.Query != "" {
		conditions = append(conditions, fmt.Sprintf("search_vector @@ plainto_tsquery('simple', $%d)", argIdx))
		args = append(args, params.Query)
		argIdx++
	}
	if params.City != "" {
		conditions = append(conditions, fmt.Sprintf("location::jsonb->>'city' = $%d", argIdx))
		args = append(args, params.City)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs WHERE %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count search: %w", err)
	}

	query := fmt.Sprintf(`SELECT id, hr_user_id, title, description, requirements::text, salary_range::text, location::text, status, to_char(expire_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM jobs WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("search jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, *j)
	}
	if jobs == nil {
		jobs = []Job{}
	}
	return jobs, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Job, error) {
	query := `SELECT id, hr_user_id, title, description, requirements::text, salary_range::text, location::text, status, to_char(expire_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM jobs WHERE id = $1`
	j, err := scanJob(r.pool.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("job not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	return j, nil
}

func (r *Repository) Create(ctx context.Context, hrUserID uuid.UUID, title, description, requirements, salaryRange, location string) (*Job, error) {
	query := `
		INSERT INTO jobs (hr_user_id, title, description, requirements, salary_range, location)
		VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6::jsonb)
		RETURNING id, hr_user_id, title, description, requirements::text, salary_range::text, location::text, status, to_char(expire_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`
	j, err := scanJob(r.pool.QueryRow(ctx, query, hrUserID, title, description, requirements, salaryRange, location))
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}
	return j, nil
}

func (r *Repository) Update(ctx context.Context, id, hrUserID uuid.UUID, title, description, requirements, salaryRange, location string) error {
	query := `UPDATE jobs SET title=$3, description=$4, requirements=$5::jsonb, salary_range=$6::jsonb, location=$7::jsonb, updated_at=NOW() WHERE id=$1 AND hr_user_id=$2`
	result, err := r.pool.Exec(ctx, query, id, hrUserID, title, description, requirements, salaryRange, location)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or not owned by you")
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id, hrUserID uuid.UUID, status string) error {
	query := `UPDATE jobs SET status=$3, updated_at=NOW() WHERE id=$1 AND hr_user_id=$2`
	result, err := r.pool.Exec(ctx, query, id, hrUserID, status)
	if err != nil {
		return fmt.Errorf("update job status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or not owned by you")
	}
	return nil
}

func (r *Repository) ListByHR(ctx context.Context, hrUserID uuid.UUID) ([]Job, error) {
	query := `SELECT id, hr_user_id, title, description, requirements::text, salary_range::text, location::text, status, to_char(expire_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') FROM jobs WHERE hr_user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, hrUserID)
	if err != nil {
		return nil, fmt.Errorf("list hr jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, *j)
	}
	if jobs == nil {
		jobs = []Job{}
	}
	return jobs, nil
}
