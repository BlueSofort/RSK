package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Post 文章/公告表
type Post struct {
	ID          uint           `gorm:"primarykey" json:"id"`                               // 主键
	Slug        string         `gorm:"uniqueIndex:idx_post_slug_deleted;not null" json:"slug"` // 唯一标识（复合索引）
	Type        string         `gorm:"not null;index" json:"type"`                         // 类型（blog/notice）
	TitleJSON   JSON           `gorm:"type:json;not null" json:"title"`                    // 多语言标题
	SummaryJSON JSON           `gorm:"type:json" json:"summary"`                           // 多语言摘要
	ContentJSON JSON           `gorm:"type:json" json:"content"`                           // 多语言内容
	Thumbnail   string         `json:"thumbnail"`                                          // 缩略图
	IsPublished bool           `gorm:"default:false;index" json:"is_published"`            // 是否发布
	PublishedAt *time.Time     `gorm:"index" json:"published_at"`                          // 发布时间
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`                            // 创建时间
	DeletedAt   gorm.DeletedAt `gorm:"uniqueIndex:idx_post_slug_deleted;index" json:"-"`   // 软删除时间（加入复合索引）
}

// MarshalJSON 自定义 JSON 序列化，精简日期显示
func (p Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	publishedAt := ""
	if p.PublishedAt != nil {
		publishedAt = p.PublishedAt.Format("2006-01-02")
	} else {
		publishedAt = p.CreatedAt.Format("2006-01-02")
	}

	return json.Marshal(&struct {
		Alias
		PublishedAtStr string `json:"published_at"`
	}{
		Alias:          (Alias)(p),
		PublishedAtStr: publishedAt,
	})
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}
