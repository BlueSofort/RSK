package repository

import (
	"github.com/dujiao-next/internal/models"

	"gorm.io/gorm"
)

// CommentRepository 评论数据访问接口
type CommentRepository interface {
	List(filter CommentListFilter) ([]models.Comment, int64, error)
	GetByID(id uint) (*models.Comment, error)
	Create(comment *models.Comment) error
	Delete(id uint) error
	GetReplies(parentIDs []uint) ([]models.Comment, error)
	CountByParentID(parentID uint) (int64, error)
}

// GormCommentRepository GORM 实现
type GormCommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *gorm.DB) *GormCommentRepository {
	return &GormCommentRepository{db: db}
}

// List 评论列表（按文章分页）
func (r *GormCommentRepository) List(filter CommentListFilter) ([]models.Comment, int64, error) {
	query := r.db.Model(&models.Comment{}).Preload("User")

	if filter.PostID > 0 {
		query = query.Where("post_id = ?", filter.PostID)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = applyPagination(query, filter.Page, filter.PageSize)

	var comments []models.Comment
	if err := query.Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// GetByID 根据 ID 获取评论
func (r *GormCommentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.Preload("User").First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// Create 创建评论
func (r *GormCommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// Delete 物理删除评论
func (r *GormCommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

// GetReplies 批量获取指定父评论的回复
func (r *GormCommentRepository) GetReplies(parentIDs []uint) ([]models.Comment, error) {
	if len(parentIDs) == 0 {
		return nil, nil
	}
	var replies []models.Comment
	if err := r.db.Preload("User").
		Where("parent_id IN ?", parentIDs).
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// CountByParentID 统计某评论的回复数
func (r *GormCommentRepository) CountByParentID(parentID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Comment{}).Where("parent_id = ?", parentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
