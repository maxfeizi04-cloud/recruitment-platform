package resume

import (
	"net/http"
	"path/filepath"
	"strings"

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

// @Summary      获取我的简历列表
// @Tags         简历
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes [get]
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

// @Summary      创建简历
// @Tags         简历
// @Accept       json
// @Produce      json
// @Param        body body resumeReq true "简历标题和内容"
// @Success      200 {object} Resume
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes [post]
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

// @Summary      更新简历
// @Tags         简历
// @Accept       json
// @Produce      json
// @Param        id path string true "简历ID"
// @Param        body body resumeReq true "简历标题和内容"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes/{id} [put]
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

// @Summary      删除简历
// @Tags         简历
// @Accept       json
// @Produce      json
// @Param        id path string true "简历ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "简历已删除"})
}

// @Summary      设为默认简历
// @Tags         简历
// @Accept       json
// @Produce      json
// @Param        id path string true "简历ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes/{id}/set-default [post]
func (h *Handler) SetDefault(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	if err := h.svc.SetDefault(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已设为默认简历"})
}

// @Summary      上传简历附件
// @Tags         简历
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "简历ID"
// @Param        file formData file true "附件文件（PDF/DOC/DOCX/JPG/PNG，最大20MB）"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /resumes/{id}/upload-attachment [post]
func (h *Handler) UploadAttachment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 验证文件类型
	allowedExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true,
		".jpg": true, ".jpeg": true, ".png": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的文件格式，仅支持 PDF、DOC、DOCX、JPG、PNG",
		})
		return
	}

	// 验证文件大小 (最大 20MB)
	const maxSize = 20 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件大小不能超过 20MB",
		})
		return
	}

	url, err := h.svc.UploadAttachment(c.Request.Context(), id, userID.(string), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
