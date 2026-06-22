package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(public, protected *gin.RouterGroup) {
	apps := protected.Group("/applications")
	{
		apps.POST("", h.Apply)
		apps.GET("/my", h.ListByCandidate)
		apps.GET("/received", h.ListByHR)
		apps.PATCH("/:id/status", h.UpdateStatus)
	}
}

type applyReq struct {
	JobID    string `json:"job_id" binding:"required"`
	ResumeID string `json:"resume_id" binding:"required"`
}

func (h *Handler) Apply(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req applyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择职位和简历"})
		return
	}
	app, err := h.svc.Apply(c.Request.Context(), userID.(string), req.JobID, req.ResumeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, app)
}

func (h *Handler) ListByCandidate(c *gin.Context) {
	userID, _ := c.Get("user_id")
	apps, err := h.svc.ListByCandidate(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"applications": apps})
}

func (h *Handler) ListByHR(c *gin.Context) {
	userID, _ := c.Get("user_id")
	apps, err := h.svc.ListByHR(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"applications": apps})
}

type updateStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req updateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态不能为空"})
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), userID.(string), c.Param("id"), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "状态已更新"})
}
