package public

import (
	"path/filepath"

	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"

	"github.com/gin-gonic/gin"
)

// UploadAvatar 用户上传头像
func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		shared.RespondError(c, response.CodeBadRequest, "error.upload_file_required", nil)
		return
	}

	// 校验文件类型（仅图片）
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		shared.RespondError(c, response.CodeBadRequest, "error.upload_file_type_invalid", nil)
		return
	}

	// 校验大小（最大 2MB）
	if file.Size > 2*1024*1024 {
		shared.RespondError(c, response.CodeBadRequest, "error.upload_file_too_large", nil)
		return
	}

	// 使用 UploadService 保存文件
	avatarURL, err := h.UploadService.SaveFile(file, "avatar")
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.upload_failed", err)
		return
	}

	// 更新用户头像字段
	_, err = h.UserAuthService.UpdateProfile(userID, nil, nil, &avatarURL)
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.profile_update_failed", err)
		return
	}

	response.Success(c, gin.H{"avatar": avatarURL})
}
