package models

import (
	"time"
)

// Comment 文章评论表
type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id"`                     // 主键
	UserID    uint      `gorm:"not null;index:idx_comments_user_id" json:"user_id"`     // 评论用户ID
	PostID    uint      `gorm:"not null;default:0;index:idx_comments_post_id" json:"post_id"` // 文章ID
	ParentID  uint      `gorm:"not null;default:0;index:idx_comments_parent_id" json:"parent_id"` // 父评论ID（0=顶层）
	Content   string    `gorm:"type:text;not null" json:"content"`        // 评论内容（纯文本，≤100字）
	Status    string    `gorm:"type:varchar(20);not null;default:'approved';index:idx_comments_status" json:"status"` // 状态（approved/rejected）
	CreatedAt time.Time `gorm:"index" json:"created_at"`                  // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                               // 更新时间

	// 关联
	User   User     `gorm:"foreignKey:UserID" json:"user,omitempty"` // 评论用户
	Parent *Comment `gorm:"foreignKey:ParentID" json:"-"`            // 父评论（不序列化）
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
