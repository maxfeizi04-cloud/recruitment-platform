package user

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/me", h.GetProfile)
		users.PUT("/me", h.UpdateProfile)
	}
	hrs := r.Group("/hrs")
	{
		hrs.POST("/certify", h.SubmitCertification)
		hrs.GET("/certification", h.GetCertification)
	}
}

// @Summary      获取个人资料
// @Tags         用户
// @Accept       json
// @Produce      json
// @Success      200 {object} Profile
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.svc.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

type updateProfileReq struct {
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

// @Summary      更新个人资料
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body body updateProfileReq true "姓名和头像"
// @Success      200 {object} Profile
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "姓名不能为空"})
		return
	}
	profile, err := h.svc.UpdateProfile(c.Request.Context(), userID.(string), req.Name, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

type certifyReq struct {
	CompanyName string `json:"company_name" binding:"required"`
	Position    string `json:"position" binding:"required"`
}

// @Summary      提交HR认证
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body body certifyReq true "公司名称和职位"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /hrs/certify [post]
func (h *Handler) SubmitCertification(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req certifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公司名称和职位不能为空"})
		return
	}
	if err := h.svc.SubmitCertification(c.Request.Context(), userID.(string), req.CompanyName, req.Position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "认证信息已提交"})
}

// @Summary      获取HR认证信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /hrs/certification [get]
func (h *Handler) GetCertification(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cert, err := h.svc.GetCertification(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if cert == nil {
		c.JSON(http.StatusOK, gin.H{"certification": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"certification": cert})
}
