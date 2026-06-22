package recommend

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobMatch struct {
	JobID      uuid.UUID `json:"job_id"`
	Title      string    `json:"title"`
	Similarity float64   `json:"similarity"`
	Salary     string    `json:"salary_range"`
	Location   string    `json:"location"`
	CreatedAt  string    `json:"created_at"`
}

type CandidateMatch struct {
	UserID     uuid.UUID `json:"user_id"`
	UserName   string    `json:"user_name"`
	ResumeID   uuid.UUID `json:"resume_id"`
	ResumeTitle string   `json:"resume_title"`
	Similarity float64   `json:"similarity"`
	Phone      string    `json:"phone"`
}

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// RecommendJobs finds top 20 jobs matching the candidate's default resume.
func (s *Service) RecommendJobs(ctx context.Context, userID uuid.UUID) ([]JobMatch, error) {
	// Get the user's default resume skills
	resumeSkills, err := s.getResumeSkills(ctx, userID)
	if err != nil || len(resumeSkills) == 0 {
		return s.fallbackJobs(ctx)
	}

	// Get all active jobs with requirements
	rows, err := s.pool.Query(ctx, `
		SELECT id, title, requirements::text, salary_range::text, location::text,
			to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM jobs WHERE status = 'active' ORDER BY created_at DESC LIMIT 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []JobMatch
	for rows.Next() {
		var id uuid.UUID
		var title, reqStr, salary, location, createdAt string
		if err := rows.Scan(&id, &title, &reqStr, &salary, &location, &createdAt); err != nil {
			continue
		}

		jobSkills := parseSkills(reqStr)
		if len(jobSkills) == 0 {
			continue
		}

		sim := jaccardSimilarity(resumeSkills, jobSkills)
		if sim > 0 {
			// Time decay: newer jobs get boost
			decay := timeDecay(createdAt)
			matches = append(matches, JobMatch{
				JobID:      id,
				Title:      title,
				Similarity: math.Round((sim*0.7+decay*0.3)*100) / 100,
				Salary:     salary,
				Location:   location,
				CreatedAt:  createdAt,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	if len(matches) > 20 {
		matches = matches[:20]
	}
	if len(matches) == 0 {
		return s.fallbackJobs(ctx)
	}
	return matches, nil
}

// RecommendCandidates finds top 20 candidates matching a job's requirements.
func (s *Service) RecommendCandidates(ctx context.Context, jobID uuid.UUID) ([]CandidateMatch, error) {
	// Get job requirements
	var reqStr string
	err := s.pool.QueryRow(ctx, `SELECT requirements::text FROM jobs WHERE id = $1`, jobID).Scan(&reqStr)
	if err != nil {
		return nil, err
	}

	jobSkills := parseSkills(reqStr)
	if len(jobSkills) == 0 {
		return nil, nil
	}

	// Get all default resumes
	rows, err := s.pool.Query(ctx, `
		SELECT r.id, r.user_id, r.title, r.content::text, u.name, u.phone
		FROM resumes r JOIN users u ON r.user_id = u.id
		WHERE r.is_default = true AND u.role = 'candidate'
		LIMIT 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []CandidateMatch
	for rows.Next() {
		var resumeID, userID uuid.UUID
		var resumeTitle, contentStr, userName, phone string
		if err := rows.Scan(&resumeID, &userID, &resumeTitle, &contentStr, &userName, &phone); err != nil {
			continue
		}

		resumeSkills := parseSkills(contentStr)
		if len(resumeSkills) == 0 {
			continue
		}

		sim := jaccardSimilarity(jobSkills, resumeSkills)
		if sim > 0 {
			matches = append(matches, CandidateMatch{
				UserID:      userID,
				UserName:    userName,
				ResumeID:    resumeID,
				ResumeTitle: resumeTitle,
				Similarity:  math.Round(sim*100) / 100,
				Phone:       maskPhone(phone),
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	if len(matches) > 20 {
		matches = matches[:20]
	}
	return matches, nil
}

func (s *Service) getResumeSkills(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var contentStr string
	err := s.pool.QueryRow(ctx,
		`SELECT content::text FROM resumes WHERE user_id = $1 AND is_default = true LIMIT 1`,
		userID,
	).Scan(&contentStr)
	if err != nil {
		// Try any resume
		err = s.pool.QueryRow(ctx,
			`SELECT content::text FROM resumes WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`,
			userID,
		).Scan(&contentStr)
		if err != nil {
			return nil, err
		}
	}
	return parseSkills(contentStr), nil
}

func (s *Service) fallbackJobs(ctx context.Context) ([]JobMatch, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, title, requirements::text, salary_range::text, location::text,
			to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM jobs WHERE status = 'active' ORDER BY created_at DESC LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []JobMatch
	for rows.Next() {
		var j JobMatch
		var reqStr string
		if err := rows.Scan(&j.JobID, &j.Title, &reqStr, &j.Salary, &j.Location, &j.CreatedAt); err != nil {
			continue
		}
		j.Similarity = 0
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// parseSkills extracts skill tags from a JSON resume/job content.
// Expects: {"skills": ["Go", "React", ...]} or just an array ["Go", "React"]
func parseSkills(jsonStr string) []string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// Try array format
		var arr []string
		if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
			return nil
		}
		return arr
	}

	if skills, ok := data["skills"]; ok {
		if arr, ok := skills.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, s := range arr {
				if str, ok := s.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

// jaccardSimilarity computes |A ∩ B| / |A ∪ B|
func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	setA := make(map[string]bool, len(a))
	for _, s := range a {
		setA[s] = true
	}

	intersection := 0
	setB := make(map[string]bool, len(b))
	for _, s := range b {
		setB[s] = true
		if setA[s] {
			intersection++
		}
	}

	// Union = all unique from both
	union := make(map[string]bool)
	for s := range setA {
		union[s] = true
	}
	for s := range setB {
		union[s] = true
	}

	if len(union) == 0 {
		return 0
	}
	return float64(intersection) / float64(len(union))
}

// timeDecay returns a decay factor based on how recent the job is.
// 0-7 days: 1.0, 7-30 days: 0.8, 30-90 days: 0.5, >90 days: 0.2
func timeDecay(createdAt string) float64 {
	t, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
	if err != nil {
		return 0.5
	}
	days := time.Since(t).Hours() / 24
	switch {
	case days <= 7:
		return 1.0
	case days <= 30:
		return 0.8
	case days <= 90:
		return 0.5
	default:
		return 0.2
	}
}

func maskPhone(phone string) string {
	if len(phone) == 11 {
		return phone[:3] + "****" + phone[7:]
	}
	return phone
}
