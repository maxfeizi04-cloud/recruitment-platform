package resume

import (
	"context"
	"fmt"
	"mime/multipart"

	"recruitment-platform/internal/pkg/cos"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	uploader cos.Uploader
}

func NewService(repo *Repository, uploader cos.Uploader) *Service {
	return &Service{repo: repo, uploader: uploader}
}

func (s *Service) List(ctx context.Context, userIDStr string) ([]Resume, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) Create(ctx context.Context, userIDStr, title, content string) (*Resume, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		content = "{}"
	}
	return s.repo.Create(ctx, userID, title, content)
}

func (s *Service) Update(ctx context.Context, idStr, userIDStr, title, content string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid resume id: %w", err)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	if title == "" {
		return fmt.Errorf("title is required")
	}
	return s.repo.Update(ctx, id, userID, title, content)
}

func (s *Service) Delete(ctx context.Context, idStr, userIDStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid resume id: %w", err)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.Delete(ctx, id, userID)
}

func (s *Service) SetDefault(ctx context.Context, idStr, userIDStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid resume id: %w", err)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.SetDefault(ctx, id, userID)
}

func (s *Service) UploadAttachment(ctx context.Context, idStr, userIDStr string, file *multipart.FileHeader) (string, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return "", fmt.Errorf("invalid resume id: %w", err)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user id: %w", err)
	}

	// Verify resume ownership
	_, err = s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return "", err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	// Generate COS key and upload
	key := cos.GenerateKey(userIDStr, file.Filename)
	url, err := s.uploader.Upload(ctx, key, src, file.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("upload to cos: %w", err)
	}

	// Add URL to resume attachment list
	if err := s.repo.AddAttachment(ctx, id, userID, url); err != nil {
		return "", err
	}

	return url, nil
}
