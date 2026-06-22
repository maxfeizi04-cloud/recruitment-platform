package recommend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(protected *gin.RouterGroup) {
	rec := protected.Group("/recommend")
	{
		rec.GET("/jobs", h.RecommendJobs)
		rec.GET("/candidates", h.RecommendCandidates)
	}
}

// @Summary      推荐职位
// @Tags         推荐
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /recommend/jobs [get]
func (h *Handler) RecommendJobs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	jobs, err := h.svc.RecommendJobs(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recommendations": jobs})
}

// @Summary      推荐候选人
// @Tags         推荐
// @Accept       json
// @Produce      json
// @Param        job_id query string true "职位ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /recommend/candidates [get]
func (h *Handler) RecommendCandidates(c *gin.Context) {
	jobID := c.Query("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}

	jid, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	candidates, err := h.svc.RecommendCandidates(c.Request.Context(), jid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if candidates == nil {
		candidates = []CandidateMatch{}
	}
	c.JSON(http.StatusOK, gin.H{"recommendations": candidates})
}
