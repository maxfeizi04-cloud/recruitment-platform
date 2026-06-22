package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(public, protected *gin.RouterGroup) {
	// 公开路由
	jobs := public.Group("/jobs")
	{
		jobs.GET("", h.List)
		jobs.GET("/search", h.Search)
		jobs.GET("/:id", h.GetByID)
	}

	// HR 路由（需要认证+HR角色）
	hr := protected.Group("/jobs")
	{
		hr.POST("", h.Create)
		hr.PUT("/:id", h.Update)
		hr.PATCH("/:id/status", h.UpdateStatus)
		hr.GET("/my", h.ListByHR)
	}
}

func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	jobs, total, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs, "total": total})
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	city := c.Query("city")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	jobs, total, err := h.svc.Search(c.Request.Context(), q, city, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs, "total": total})
}

func (h *Handler) GetByID(c *gin.Context) {
	job, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
		return
	}
	c.JSON(http.StatusOK, job)
}

type jobReq struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Requirements string `json:"requirements"`
	SalaryRange  string `json:"salary_range"`
	Location     string `json:"location"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "hr" {
		c.JSON(http.StatusForbidden, gin.H{"error": "仅HR可发布职位"})
		return
	}
	var req jobReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "职位标题不能为空"})
		return
	}
	job, err := h.svc.Create(c.Request.Context(), userID.(string), req.Title, req.Description, req.Requirements, req.SalaryRange, req.Location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *Handler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req jobReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "职位标题不能为空"})
		return
	}
	if err := h.svc.Update(c.Request.Context(), c.Param("id"), userID.(string), req.Title, req.Description, req.Requirements, req.SalaryRange, req.Location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "职位已更新"})
}

type statusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态值不能为空"})
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), c.Param("id"), userID.(string), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "状态已更新"})
}

func (h *Handler) ListByHR(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobs, err := h.svc.ListByHR(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
