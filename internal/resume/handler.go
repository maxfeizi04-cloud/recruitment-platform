package resume

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
	resumes := r.Group("/resumes")
	{
		resumes.GET("", h.List)
		resumes.POST("", h.Create)
		resumes.PUT("/:id", h.Update)
		resumes.DELETE("/:id", h.Delete)
		resumes.POST("/:id/set-default", h.SetDefault)
		resumes.POST("/:id/upload-attachment", h.UploadAttachment)
	}
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resumes, err := h.svc.List(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"resumes": resumes})
}

type resumeReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req resumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "简历标题不能为空"})
		return
	}
	resume, err := h.svc.Create(c.Request.Context(), userID.(string), req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resume)
}

func (h *Handler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	var req resumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "简历标题不能为空"})
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, userID.(string), req.Title, req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "简历已更新"})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "简历已删除"})
}

func (h *Handler) SetDefault(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	if err := h.svc.SetDefault(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已设为默认简历"})
}

func (h *Handler) UploadAttachment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	url, err := h.svc.UploadAttachment(c.Request.Context(), id, userID.(string), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
