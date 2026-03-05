package admin

import (
	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"

	"github.com/gin-gonic/gin"
)

// ====================  文件上传  ====================

// UploadFile 文件上传
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		shared.RespondError(c, response.CodeBadRequest, "error.file_missing", nil)
		return
	}
	scene := c.DefaultPostForm("scene", "common")

	// 保存文件
	url, err := h.UploadService.SaveFile(file, scene)
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.upload_failed", err)
		return
	}

	response.Success(c, gin.H{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
