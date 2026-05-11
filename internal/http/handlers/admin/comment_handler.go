package admin

import (
	"errors"
	"strconv"

	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"
	"github.com/dujiao-next/internal/service"

	"github.com/gin-gonic/gin"
)

// ────────── 管理端评论接口 ──────────

// GetAdminComments 管理员获取评论列表
func (h *Handler) GetAdminComments(c *gin.Context) {
	postID, _ := strconv.Atoi(c.DefaultQuery("post_id", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	page, pageSize = shared.NormalizePagination(page, pageSize)

	comments, total, err := h.CommentService.ListAdmin(postID, page, pageSize)
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.comment_fetch_failed", err)
		return
	}

	pagination := response.BuildPagination(page, pageSize, total)
	response.SuccessWithPage(c, comments, pagination)
}

// DeleteAdminComment 管理员删除评论（物理删除）
func (h *Handler) DeleteAdminComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.RespondError(c, response.CodeBadRequest, "error.bad_request", nil)
		return
	}

	if err := h.CommentService.DeleteByAdmin(uint(commentID)); err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			shared.RespondError(c, response.CodeNotFound, "error.comment_not_found", nil)
			return
		}
		shared.RespondError(c, response.CodeInternal, "error.comment_delete_failed", err)
		return
	}

	response.Success(c, nil)
}
