package interview

import (
	"encoding/json"
	"net/http"

	"recruitment-platform/internal/pkg/maps"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc        *Service
	mapsClient *maps.Client
}

func NewHandler(svc *Service, mapsClient *maps.Client) *Handler {
	return &Handler{svc: svc, mapsClient: mapsClient}
}

func (h *Handler) RegisterRoutes(public, protected *gin.RouterGroup) {
	interviews := protected.Group("/interviews")
	{
		interviews.POST("", h.Create)
		interviews.GET("/my", h.ListByUser)
		interviews.PATCH("/:id/status", h.UpdateStatus)
		interviews.GET("/:id/navigate", h.Navigate)
	}
	public.GET("/maps/place-search", h.PlaceSearch)
}

type interviewReq struct {
	ApplicationID string `json:"application_id" binding:"required"`
	ScheduledAt   string `json:"scheduled_at" binding:"required"`
	Address       string `json:"company_address" binding:"required"`
	ContactName   string `json:"contact_name"`
	ContactPhone  string `json:"contact_phone"`
	Notes         string `json:"notes"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "hr" {
		c.JSON(http.StatusForbidden, gin.H{"error": "仅HR可发起面试邀约"})
		return
	}
	var req interviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整信息"})
		return
	}
	inv, err := h.svc.Create(c.Request.Context(), userID.(string), req.ApplicationID, req.ScheduledAt, req.Address, req.ContactName, req.ContactPhone, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *Handler) ListByUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	invs, err := h.svc.ListByUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invitations": invs})
}

type updateInvStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req updateInvStatusReq
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

func (h *Handler) Navigate(c *gin.Context) {
	inv, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀约不存在"})
		return
	}
	var addr struct {
		Lng       float64 `json:"lng"`
		Lat       float64 `json:"lat"`
		Formatted string  `json:"formatted"`
	}
	json.Unmarshal([]byte(inv.CompanyAddress), &addr)
	navURL := maps.GenerateNavigationURL(addr.Lat, addr.Lng, addr.Formatted)
	c.JSON(http.StatusOK, gin.H{"navigation_url": navURL})
}

func (h *Handler) PlaceSearch(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	results, err := h.mapsClient.PlaceSearch(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"suggestions": results})
}
