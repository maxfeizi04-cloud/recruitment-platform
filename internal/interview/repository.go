package interview

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invitation struct {
	ID               uuid.UUID `json:"id"`
	JobApplicationID uuid.UUID `json:"job_application_id"`
	HRUserID         uuid.UUID `json:"hr_user_id"`
	CandidateUserID  uuid.UUID `json:"candidate_user_id"`
	ScheduledAt      string    `json:"scheduled_at"`
	CompanyAddress   string    `json:"company_address"`
	ContactName      string    `json:"contact_name"`
	ContactPhone     string    `json:"contact_phone"`
	Notes            string    `json:"notes"`
	AttachmentURLs   []string  `json:"attachment_urls"`
	Status           string    `json:"status"`
	JobTitle         string    `json:"job_title,omitempty"`
	CandidateName    string    `json:"candidate_name,omitempty"`
	HRName           string    `json:"hr_name,omitempty"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, hrUserID, candidateUserID, appID uuid.UUID, scheduledAt, address, contactName, contactPhone, notes string) (*Invitation, error) {
	query := `
		INSERT INTO interview_invitations (job_application_id, hr_user_id, candidate_user_id, scheduled_at, company_address, contact_name, contact_phone, notes)
		VALUES ($1, $2, $3, $4::timestamp, $5::jsonb, $6, $7, $8)
		RETURNING id, job_application_id, hr_user_id, candidate_user_id,
			to_char(scheduled_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			company_address::text, contact_name, contact_phone, notes, attachment_urls, status,
			to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`
	var inv Invitation
	var attachmentURLs []string
	err := r.pool.QueryRow(ctx, query, appID, hrUserID, candidateUserID, scheduledAt, address, contactName, contactPhone, notes).Scan(
		&inv.ID, &inv.JobApplicationID, &inv.HRUserID, &inv.CandidateUserID,
		&inv.ScheduledAt, &inv.CompanyAddress, &inv.ContactName, &inv.ContactPhone, &inv.Notes, &attachmentURLs, &inv.Status,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create invitation: %w", err)
	}
	if attachmentURLs == nil {
		attachmentURLs = []string{}
	}
	inv.AttachmentURLs = attachmentURLs
	return &inv, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Invitation, error) {
	query := `
		SELECT i.id, i.job_application_id, i.hr_user_id, i.candidate_user_id,
			COALESCE(to_char(i.scheduled_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
			i.company_address::text, i.contact_name, i.contact_phone, i.notes, i.attachment_urls, i.status,
			j.title AS job_title, u_hr.name AS hr_name, u_cand.name AS candidate_name,
			to_char(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(i.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM interview_invitations i
		JOIN job_applications a ON i.job_application_id = a.id
		JOIN jobs j ON a.job_id = j.id
		JOIN users u_hr ON i.hr_user_id = u_hr.id
		JOIN users u_cand ON i.candidate_user_id = u_cand.id
		WHERE i.hr_user_id = $1 OR i.candidate_user_id = $1
		ORDER BY i.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list invitations: %w", err)
	}
	defer rows.Close()

	var invs []Invitation
	for rows.Next() {
		var inv Invitation
		var attachmentURLs []string
		if err := rows.Scan(&inv.ID, &inv.JobApplicationID, &inv.HRUserID, &inv.CandidateUserID,
			&inv.ScheduledAt, &inv.CompanyAddress, &inv.ContactName, &inv.ContactPhone, &inv.Notes, &attachmentURLs, &inv.Status,
			&inv.JobTitle, &inv.HRName, &inv.CandidateName, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		if attachmentURLs == nil {
			attachmentURLs = []string{}
		}
		inv.AttachmentURLs = attachmentURLs
		invs = append(invs, inv)
	}
	if invs == nil {
		invs = []Invitation{}
	}
	return invs, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	query := `
		SELECT i.id, i.job_application_id, i.hr_user_id, i.candidate_user_id,
			COALESCE(to_char(i.scheduled_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
			i.company_address::text, i.contact_name, i.contact_phone, i.notes, i.attachment_urls, i.status,
			j.title, u_hr.name, u_cand.name,
			to_char(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			to_char(i.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM interview_invitations i
		JOIN job_applications a ON i.job_application_id = a.id
		JOIN jobs j ON a.job_id = j.id
		JOIN users u_hr ON i.hr_user_id = u_hr.id
		JOIN users u_cand ON i.candidate_user_id = u_cand.id
		WHERE i.id = $1
	`
	var inv Invitation
	var attachmentURLs []string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.JobApplicationID, &inv.HRUserID, &inv.CandidateUserID,
		&inv.ScheduledAt, &inv.CompanyAddress, &inv.ContactName, &inv.ContactPhone, &inv.Notes, &attachmentURLs, &inv.Status,
		&inv.JobTitle, &inv.HRName, &inv.CandidateName, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("invitation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get invitation: %w", err)
	}
	if attachmentURLs == nil {
		attachmentURLs = []string{}
	}
	inv.AttachmentURLs = attachmentURLs
	return &inv, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status string) error {
	query := `UPDATE interview_invitations SET status = $3, updated_at = NOW() WHERE id = $1 AND (hr_user_id = $2 OR candidate_user_id = $2)`
	result, err := r.pool.Exec(ctx, query, id, userID, status)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invitation not found or not authorized")
	}
	return nil
}
