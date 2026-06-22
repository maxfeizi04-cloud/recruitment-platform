package auth

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
	auth := r.Group("/auth")
	{
		auth.POST("/send-code", h.SendCode)
		auth.POST("/login", h.Login)
	}
}

type sendCodeReq struct {
	Phone string `json:"phone" binding:"required,len=11"`
}

func (h *Handler) SendCode(c *gin.Context) {
	var req sendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	if err := h.svc.SendVerificationCode(c.Request.Context(), req.Phone); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

type loginReq struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Code  string `json:"code" binding:"required,len=6"`
	Role  string `json:"role" binding:"required,oneof=candidate hr"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写正确的手机号、验证码和角色"})
		return
	}

	result, err := h.svc.Login(c.Request.Context(), req.Phone, req.Code, req.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}
