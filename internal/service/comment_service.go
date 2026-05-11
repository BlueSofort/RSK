package service

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dujiao-next/internal/cache"
	"github.com/dujiao-next/internal/config"
	"github.com/dujiao-next/internal/models"
	"github.com/dujiao-next/internal/repository"
)

//go:embed sensitive_words.txt
var embeddedSensitiveWords string

// CommentService 评论业务服务
type CommentService struct {
	repo     repository.CommentRepository
	postRepo repository.PostRepository
	cfg      *config.Config
	wordTrie *sensitiveTrie // DFA 敏感词字典树
}

// NewCommentService 创建评论服务
func NewCommentService(
	repo repository.CommentRepository,
	postRepo repository.PostRepository,
	cfg *config.Config,
) *CommentService {
	s := &CommentService{repo: repo, postRepo: postRepo, cfg: cfg}
	if cfg.Security.SensitiveWords.Enabled {
		s.wordTrie = loadSensitiveWords(cfg.Security.SensitiveWords.Dict)
	}
	return s
}

// ────────── DFA 敏感词过滤 ──────────

type sensitiveTrie struct {
	root map[rune]*sensitiveNode
}

type sensitiveNode struct {
	children map[rune]*sensitiveNode
	isEnd    bool
}

func loadSensitiveWords(path string) *sensitiveTrie {
	trie := &sensitiveTrie{root: make(map[rune]*sensitiveNode)}

	var scanner *bufio.Scanner
	if strings.TrimSpace(embeddedSensitiveWords) != "" {
		// 优先使用编译进二进制的词库
		scanner = bufio.NewScanner(strings.NewReader(embeddedSensitiveWords))
	} else if path != "" {
		// 回退到配置文件指定的外部文件
		f, err := os.Open(path)
		if err != nil {
			return trie
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	} else {
		return trie
	}
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		trie.insert(word)
	}
	return trie
}

func (t *sensitiveTrie) insert(word string) {
	if word == "" {
		return
	}
	runes := []rune(word)

	// 初始化根节点
	root := t.root
	first := runes[0]
	if _, ok := root[first]; !ok {
		root[first] = &sensitiveNode{children: make(map[rune]*sensitiveNode)}
	}

	node := root[first]
	for i := 1; i < len(runes); i++ {
		r := runes[i]
		if next, ok := node.children[r]; ok {
			node = next
		} else {
			newNode := &sensitiveNode{children: make(map[rune]*sensitiveNode)}
			node.children[r] = newNode
			node = newNode
		}
	}
	// 最后一个节点标记为词尾
	if node != nil {
		node.isEnd = true
	}
}

// checkSensitive DFA 扫描文本，返回第一个命中的敏感词
func (t *sensitiveTrie) check(content string) (string, bool) {
	if t == nil || len(t.root) == 0 {
		return "", false
	}
	runes := []rune(content)
	n := len(runes)
	for i := 0; i < n; i++ {
		node, ok := t.root[runes[i]]
		if !ok {
			continue
		}
		// 最长匹配
		j := i
		for j < n && node != nil {
			if node.isEnd {
				return string(runes[i : j+1]), true
			}
			j++
			if j < n {
				node = node.children[runes[j]]
			}
		}
		// 检查是否是单字符敏感词
		if node, ok := t.root[runes[i]]; ok && node.isEnd {
			return string(runes[i]), true
		}
	}
	return "", false
}

// ────────── 评论校验 ──────────

const commentMaxChars = 100

func validateCommentContent(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ErrCommentContentEmpty
	}
	// 只允许文字和基本标点，拒绝纯空白/特殊字符
	charCount := utf8.RuneCountInString(trimmed)
	if charCount > commentMaxChars {
		return ErrCommentContentTooLong
	}
	// 校验只允许文字、数字、中文标点、英文基本标点、空格
	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		if unicode.Is(unicode.Han, r) {
			continue
		}
		// 允许基本标点
		switch r {
		case '.', ',', '!', '?', ':', ';', '-', '_', '(', ')', '，', '。', '！', '？', '：', '；', '、', '…', '~', '·':
			continue
		default:
			return ErrCommentContentInvalid
		}
	}
	return nil
}

// ────────── 业务方法 ──────────

// CreateInput 创建评论输入
type CreateCommentInput struct {
	PostID   uint   `json:"post_id"`
	ParentID uint   `json:"parent_id"`
	Content  string `json:"content"`
}

// Create 创建评论（发布即 approved）
func (s *CommentService) Create(userID uint, input CreateCommentInput) (*models.Comment, error) {
	// 1. 校验内容
	if err := validateCommentContent(input.Content); err != nil {
		return nil, err
	}

	// 2. 敏感词检测
	if s.wordTrie != nil {
		if word, hit := s.wordTrie.check(strings.TrimSpace(input.Content)); hit {
			return nil, fmt.Errorf("%w: %s", ErrCommentSensitiveWords, word)
		}
	}

	// 3. 校验文章存在
	if input.PostID == 0 {
		return nil, ErrCommentPostNotFound
	}
	post, err := s.postRepo.GetByID(fmt.Sprintf("%d", input.PostID))
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrCommentPostNotFound
	}

	// 4. 校验父评论存在（如果回复）
	if input.ParentID > 0 {
		parent, err := s.repo.GetByID(input.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.PostID != input.PostID {
			return nil, ErrCommentNotFound
		}
	}

	// 5. 频率限制
	if err := s.checkRateLimit(userID); err != nil {
		return nil, err
	}

	// 6. 创建评论
	comment := &models.Comment{
		UserID:   userID,
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Content:  strings.TrimSpace(input.Content),
		Status:   "approved",
	}
	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// GetByPostID 获取文章评论列表（含回复）
func (s *CommentService) GetByPostID(postID, page, pageSize int) ([]models.Comment, []models.Comment, int64, error) {
	parentID := uint(0)
	filter := repository.CommentListFilter{
		Page:     page,
		PageSize: pageSize,
		PostID:   uint(postID),
		ParentID: &parentID,
		Status:   "approved",
	}
	comments, total, err := s.repo.List(filter)
	if err != nil {
		return nil, nil, 0, err
	}

	// 获取回复
	if len(comments) > 0 {
		parentIDs := make([]uint, len(comments))
		for i, c := range comments {
			parentIDs[i] = c.ID
		}
		replies, err := s.repo.GetReplies(parentIDs)
		if err != nil {
			return nil, nil, 0, err
		}
		return comments, replies, total, nil
	}
	return comments, nil, total, nil
}

// DeleteByUser 用户删除自己的评论
func (s *CommentService) DeleteByUser(commentID, userID uint) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}
	if comment.UserID != userID {
		return ErrCommentNotOwner
	}

	// 有回复 → 清空内容保留壳
	count, err := s.repo.CountByParentID(commentID)
	if err != nil {
		return err
	}
	if count > 0 {
		comment.Content = "[该评论已被作者删除]"
		comment.UpdatedAt = time.Now()
		// 只更新 content 和 updated_at
		return models.DB.Model(&models.Comment{}).Where("id = ?", commentID).Updates(map[string]interface{}{
			"content":    comment.Content,
			"updated_at": comment.UpdatedAt,
		}).Error
	}

	// 无回复 → 物理删除
	return s.repo.Delete(commentID)
}

// DeleteByAdmin 管理员删除评论（物理删除，连带子回复）
func (s *CommentService) DeleteByAdmin(commentID uint) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}
	return s.repo.Delete(commentID)
}

// ListAdmin 管理员评论列表
func (s *CommentService) ListAdmin(postID, page, pageSize int) ([]models.Comment, int64, error) {
	filter := repository.CommentListFilter{
		Page:     page,
		PageSize: pageSize,
		PostID:   uint(postID),
		Status:   "approved",
	}
	return s.repo.List(filter)
}

// ────────── 频率限制 (Redis) ──────────

func (s *CommentService) checkRateLimit(userID uint) error {
	rc := cache.Client()
	if rc == nil {
		return nil // Redis 不可用则跳过
	}
	ctx := context.Background()
	prefix := strings.TrimSpace(s.cfg.Redis.Prefix)
	if prefix == "" {
		prefix = "dj"
	}
	key := fmt.Sprintf("%s:rate:comment:%d", prefix, userID)
	window := s.cfg.Security.CommentRateLimit.WindowSeconds
	maxReq := s.cfg.Security.CommentRateLimit.MaxRequests
	if window <= 0 || maxReq <= 0 {
		return nil
	}

	count, err := rc.Incr(ctx, key).Result()
	if err != nil {
		return nil // Redis 错误则跳过
	}
	if count == 1 {
		rc.Expire(ctx, key, time.Duration(window)*time.Second)
	}
	if count > int64(maxReq) {
		return ErrCommentRateLimited
	}
	return nil
}

// ────────── 默认头像 ──────────

// DefaultAvatar 根据昵称生成默认头像（首字符大写）
func DefaultAvatar(displayName string) string {
	name := strings.TrimSpace(displayName)
	if name == "" {
		return ""
	}
	first, _ := utf8.DecodeRuneInString(name)
	return strings.ToUpper(string(first))
}
