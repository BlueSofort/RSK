package public

import (
	"errors"
	"strconv"

	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"
	"github.com/dujiao-next/internal/service"

	"github.com/gin-gonic/gin"
)

// ────────── 用户端评论接口 ──────────

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

// CreateComment 用户发表评论
func (h *Handler) CreateComment(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondBindError(c, err)
		return
	}

	input := service.CreateCommentInput{
		PostID:   req.PostID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	comment, err := h.CommentService.Create(userID, input)
	if err != nil {
		respondCommentError(c, err)
		return
	}

	response.Success(c, comment)
}

// GetComments 获取文章评论列表（含回复）
func (h *Handler) GetComments(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Query("post_id"))
	if postID <= 0 {
		shared.RespondError(c, response.CodeBadRequest, "error.bad_request", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	page, pageSize = shared.NormalizePagination(page, pageSize)

	comments, replies, total, err := h.CommentService.GetByPostID(postID, page, pageSize)
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.comment_fetch_failed", err)
		return
	}

	result := map[string]interface{}{
		"list":    comments,
		"replies": replies,
	}
	pagination := response.BuildPagination(page, pageSize, total)
	response.SuccessWithPage(c, result, pagination)
}

// DeleteMyComment 用户删除自己的评论
func (h *Handler) DeleteMyComment(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.RespondError(c, response.CodeBadRequest, "error.bad_request", nil)
		return
	}

	if err := h.CommentService.DeleteByUser(uint(commentID), userID); err != nil {
		respondCommentError(c, err)
		return
	}

	response.Success(c, nil)
}

// ────────── 错误映射 ──────────

var commentErrorRules = []mappedHandlerError{
	{target: service.ErrCommentContentEmpty, code: response.CodeBadRequest, key: "error.comment_content_empty"},
	{target: service.ErrCommentContentTooLong, code: response.CodeBadRequest, key: "error.comment_content_too_long"},
	{target: service.ErrCommentContentInvalid, code: response.CodeBadRequest, key: "error.comment_content_invalid"},
	{target: service.ErrCommentSensitiveWords, code: response.CodeBadRequest, key: "error.comment_sensitive_words"},
	{target: service.ErrCommentPostNotFound, code: response.CodeNotFound, key: "error.comment_post_not_found"},
	{target: service.ErrCommentNotFound, code: response.CodeNotFound, key: "error.comment_not_found"},
	{target: service.ErrCommentNotOwner, code: response.CodeForbidden, key: "error.comment_not_owner"},
	{target: service.ErrCommentRateLimited, code: response.CodeTooManyRequests, key: "error.comment_rate_limited"},
}

func respondCommentError(c *gin.Context, err error) {
	for _, rule := range commentErrorRules {
		if errors.Is(err, rule.target) {
			shared.RespondError(c, rule.code, rule.key, nil)
			return
		}
	}
	shared.RespondError(c, response.CodeInternal, "error.comment_create_failed", err)
}
