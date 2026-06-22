package interview

import (
	"encoding/json"
	"net/http"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/maps"

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

// @Summary      发起面试邀约
// @Tags         面试
// @Accept       json
// @Produce      json
// @Param        body body interviewReq true "面试信息"
// @Success      200 {object} Invitation
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /interviews [post]
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

// @Summary      获取我的面试邀约
// @Tags         面试
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /interviews/my [get]
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

// @Summary      更新面试状态
// @Tags         面试
// @Accept       json
// @Produce      json
// @Param        id path string true "面试ID"
// @Param        body body updateInvStatusReq true "状态值（accepted/declined/reschedule/confirmed）"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /interviews/{id}/status [patch]
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

// @Summary      获取面试导航链接
// @Tags         面试
// @Accept       json
// @Produce      json
// @Param        id path string true "面试ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /interviews/{id}/navigate [get]
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

// @Summary      地点搜索
// @Tags         地图
// @Accept       json
// @Produce      json
// @Param        q query string true "搜索关键词"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /maps/place-search [get]
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
