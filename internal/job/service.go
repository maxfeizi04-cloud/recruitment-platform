package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/cache"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/search"
)

var (
	ErrEmptyTitle    = errors.New("职位标题不能为空")
	ErrEmptyLocation = errors.New("工作地址不能为空，请填写省/市")
	ErrInvalidStatus = errors.New("无效的状态值，允许: active, paused, closed")
)

type Service struct {
	repo     *Repository
	cache    *cache.Client
	esClient *search.Client
}

func NewService(repo *Repository, cache *cache.Client, esClient *search.Client) *Service {
	return &Service{repo: repo, cache: cache, esClient: esClient}
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Job, int, error) {
	if limit <= 0 || limit > 50 { limit = 20 }
	if offset < 0 { offset = 0 }

	if s.cache != nil {
		key := fmt.Sprintf("jobs:list:%d:%d", limit, offset)
		var cached struct {
			Jobs  []Job `json:"jobs"`
			Total int   `json:"total"`
		}
		if ok, _ := s.cache.Get(ctx, key, &cached); ok && len(cached.Jobs) > 0 {
			return cached.Jobs, cached.Total, nil
		}
		jobs, total, err := s.repo.List(ctx, limit, offset)
		if err == nil {
			s.cache.Set(ctx, key, struct {
				Jobs  []Job `json:"jobs"`
				Total int   `json:"total"`
			}{jobs, total}, cache.TTLMedium)
		}
		return jobs, total, err
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Search(ctx context.Context, query, city string, limit, offset int) ([]Job, int, error) {
	if limit <= 0 || limit > 50 { limit = 20 }
	if offset < 0 { offset = 0 }

	// 优先使用 ES 搜索
	if s.esClient != nil && query != "" {
		results, total, err := s.esClient.SearchJobs(ctx, query, city, limit, offset)
		if err == nil && len(results) > 0 {
			var jobs []Job
			for _, r := range results {
				id, _ := uuid.Parse(r.ID)
				job, err := s.repo.GetByID(ctx, id)
				if err == nil {
					jobs = append(jobs, *job)
				}
			}
			if len(jobs) > 0 {
				return jobs, total, nil
			}
		}
		slog.Warn("ES search fallback to PostgreSQL", "query", query, "error", err)
	}

	return s.repo.Search(ctx, SearchParams{Query: query, City: city, Limit: limit, Offset: offset})
}

func (s *Service) GetByID(ctx context.Context, idStr string) (*Job, error) {
	id, err := uuid.Parse(idStr)
	if err != nil { return nil, err }

	if s.cache != nil {
		key := fmt.Sprintf("jobs:detail:%s", idStr)
		var cached Job
		if ok, _ := s.cache.Get(ctx, key, &cached); ok {
			return &cached, nil
		}
		job, err := s.repo.GetByID(ctx, id)
		if err == nil {
			s.cache.Set(ctx, key, job, cache.TTLMedium)
		}
		return job, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, hrUserIDStr, title, description, requirements, salaryRange, location string) (*Job, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil { return nil, err }
	if title == "" { return nil, ErrEmptyTitle }
		if !validLocation(location) { return nil, ErrEmptyLocation }

	job, err := s.repo.Create(ctx, hrUserID, title, description, requirements, salaryRange, location)
	if err == nil && s.cache != nil {
		s.cache.DeletePattern(ctx, "jobs:list:*")
	}
	return job, err
}

func (s *Service) Update(ctx context.Context, idStr, hrUserIDStr, title, description, requirements, salaryRange, location string) error {
	id, err := uuid.Parse(idStr)
	if err != nil { return err }
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil { return err }
	if title == "" { return ErrEmptyTitle }
	if !validLocation(location) { return ErrEmptyLocation }

	err = s.repo.Update(ctx, id, hrUserID, title, description, requirements, salaryRange, location)
	if err == nil && s.cache != nil {
		s.cache.DeletePattern(ctx, "jobs:list:*")
		s.cache.Delete(ctx, fmt.Sprintf("jobs:detail:%s", idStr))
	}
	return err
}

func (s *Service) UpdateStatus(ctx context.Context, idStr, hrUserIDStr, status string) error {
	id, err := uuid.Parse(idStr)
	if err != nil { return err }
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil { return err }
	if status != "active" && status != "paused" && status != "closed" { return ErrInvalidStatus }

	err = s.repo.UpdateStatus(ctx, id, hrUserID, status)
	if err == nil && s.cache != nil {
		s.cache.DeletePattern(ctx, "jobs:list:*")
		s.cache.Delete(ctx, fmt.Sprintf("jobs:detail:%s", idStr))
	}
	return err
}

func (s *Service) ListByHR(ctx context.Context, hrUserIDStr string) ([]Job, error) {
	hrUserID, err := uuid.Parse(hrUserIDStr)
	if err != nil { return nil, err }
	return s.repo.ListByHR(ctx, hrUserID)
}


func validLocation(loc string) bool {
	if loc == "" || loc == "{}" { return false }
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(loc), &m); err != nil { return false }
	p, _ := m["province"].(string)
	c, _ := m["city"].(string)
	return p != "" || c != ""
}
