package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, userIDStr string) (*Profile, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.GetProfile(ctx, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, userIDStr, name, avatarURL string) (*Profile, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if err := s.repo.UpdateProfile(ctx, userID, name, avatarURL); err != nil {
		return nil, err
	}
	return s.repo.GetProfile(ctx, userID)
}

func (s *Service) SubmitCertification(ctx context.Context, userIDStr, companyName, position string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	if companyName == "" || position == "" {
		return fmt.Errorf("company name and position are required")
	}
	return s.repo.SubmitCertification(ctx, userID, companyName, position)
}

func (s *Service) GetCertification(ctx context.Context, userIDStr string) (*HRCertification, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.GetCertification(ctx, userID)
}
