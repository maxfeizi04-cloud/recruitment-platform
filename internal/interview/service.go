package interview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidStatus = errors.New("无效状态，允许: accepted, declined, reschedule, confirmed")
)

type Service struct {
	repo *Repository
	pool *pgxpool.Pool
}

func NewService(repo *Repository, pool *pgxpool.Pool) *Service {
	return &Service{repo: repo, pool: pool}
}

func (s *Service) Create(ctx context.Context, hrUserIDStr, appIDStr, scheduledAt, address, contactName, contactPhone, notes string) (*Invitation, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hr user id")
	}
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid application id")
	}
	if !json.Valid([]byte(address)) {
		return nil, fmt.Errorf("地址格式无效")
	}

	// Look up the candidate from the application
	var candidateUserID uuid.UUID
	err = s.pool.QueryRow(ctx, `SELECT user_id FROM job_applications WHERE id = $1`, appID).Scan(&candidateUserID)
	if err != nil {
		return nil, fmt.Errorf("application not found")
	}

	return s.repo.Create(ctx, hrUserID, candidateUserID, appID, scheduledAt, address, contactName, contactPhone, notes)
}

func (s *Service) ListByUser(ctx context.Context, userIDStr string) ([]Invitation, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, idStr string) (*Invitation, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateStatus(ctx context.Context, userIDStr, idStr, status string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	if status != "accepted" && status != "declined" && status != "reschedule" && status != "confirmed" {
		return ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, userID, status)
}
