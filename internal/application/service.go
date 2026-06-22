package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidStatus = errors.New("无效状态，允许: viewed, accepted, rejected")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Apply(ctx context.Context, userIDStr, jobIDStr, resumeIDStr string) (*Application, error) {
	userID, jobID, resumeID, err := parseIDs(userIDStr, jobIDStr, resumeIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, userID, jobID, resumeID)
}

func (s *Service) ListByCandidate(ctx context.Context, userIDStr string) ([]Application, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByCandidate(ctx, userID)
}

func (s *Service) ListByHR(ctx context.Context, hrUserIDStr string) ([]Application, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByHR(ctx, hrUserID)
}

func (s *Service) UpdateStatus(ctx context.Context, hrUserIDStr, appIDStr, status string) error {
	hrUserID, appID, err := parseTwoIDs(hrUserIDStr, appIDStr)
	if err != nil {
		return err
	}
	if status != "viewed" && status != "accepted" && status != "rejected" {
		return ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, appID, hrUserID, status)
}

func parseIDs(a, b, c string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	aa, err := uuid.Parse(a)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	bb, err := uuid.Parse(b)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	cc, err := uuid.Parse(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	return aa, bb, cc, nil
}

func parseTwoIDs(a, b string) (uuid.UUID, uuid.UUID, error) {
	aa, err := uuid.Parse(a)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	bb, err := uuid.Parse(b)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return aa, bb, nil
}
