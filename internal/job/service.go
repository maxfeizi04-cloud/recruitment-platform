package job

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyTitle    = errors.New("职位标题不能为空")
	ErrInvalidStatus = errors.New("无效的状态值，允许: active, paused, closed")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Job, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Search(ctx context.Context, query, city string, limit, offset int) ([]Job, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.Search(ctx, SearchParams{Query: query, City: city, Limit: limit, Offset: offset})
}

func (s *Service) GetByID(ctx context.Context, idStr string) (*Job, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, hrUserIDStr, title, description, requirements, salaryRange, location string) (*Job, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return nil, err
	}
	if title == "" {
		return nil, ErrEmptyTitle
	}
	return s.repo.Create(ctx, hrUserID, title, description, requirements, salaryRange, location)
}

func (s *Service) Update(ctx context.Context, idStr, hrUserIDStr, title, description, requirements, salaryRange, location string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return err
	}
	if title == "" {
		return ErrEmptyTitle
	}
	return s.repo.Update(ctx, id, hrUserID, title, description, requirements, salaryRange, location)
}

func (s *Service) UpdateStatus(ctx context.Context, idStr, hrUserIDStr, status string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return err
	}
	if status != "active" && status != "paused" && status != "closed" {
		return ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, hrUserID, status)
}

func (s *Service) ListByHR(ctx context.Context, hrUserIDStr string) ([]Job, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByHR(ctx, hrUserID)
}
