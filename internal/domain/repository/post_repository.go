// DB interaction logic using GORM
package repository

import (
	"fmt"
	"time"
	"context"
	"encoding/json"

	"japa/internal/domain/entity"
	"japa/internal/app/http/dto/request"

	"gorm.io/gorm"
)

// TYPES

// Repository method signatures
type PostRepositoryInterface interface {
	Create(post *entity.Post) error
	FindByID(ID string) (*entity.Post, error)
}

// PostRepository to interface with DB
type PostRepository struct {
	DB *gorm.DB
}

// METHODS

// Initialize UserRepository
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// Create post
func (pr *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	return pr.DB.WithContext(ctx).Create(post).Error
}

// Fetch posts paginated
func (pr *PostRepository) FetchPosts(ctx context.Context, limit, offset int) ([]request.PostWithAuthor, int64, error) {
	var rows []struct {
		ID          string
		Title       string
		Slug        string
		Content     string
		Excerpt     *string
		TagsRaw     string
		AccessLevel string
		CreatedAt   time.Time
		AuthorName  *string
	}

	err := pr.DB.WithContext(ctx).
		Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.slug,
			posts.content,
			posts.excerpt,
			posts.tags,
			posts.access_level,
			posts.created_at,
			users.full_name as author_name
		`).
		Joins("left join users on users.id = posts.author_id").
		Order("posts.created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert to PostWithAuthor and unmarshal Tags
	posts := make([]request.PostWithAuthor, len(rows))
	for i, r := range rows {
		posts[i] = request.PostWithAuthor{
			ID:          r.ID,
			Title:       r.Title,
			Slug:        r.Slug,
			Content:     r.Content,
			Excerpt:     r.Excerpt,
			CreatedAt:   r.CreatedAt,
			AccessLevel: r.AccessLevel,
			AuthorName:  r.AuthorName,
			TagsRaw:     r.TagsRaw,
		}

		if err := json.Unmarshal([]byte(r.TagsRaw), &posts[i].Tags); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	// Count total posts
	var totalPosts int64
	if err := pr.DB.WithContext(ctx).
		Table("posts").
		Count(&totalPosts).Error; err != nil {
		return nil, 0, err
	}

	return posts, totalPosts, nil
}

// Create post
func (pr *PostRepository) FetchPost(ctx context.Context, postID string) (*request.PostWithAuthor, error) {
	var row struct {
		ID          string
		Title       string
		Slug        string
		Content     string
		Excerpt     *string
		TagsRaw     string
		Source      *string
		AccessLevel string
		CreatedAt   time.Time
		AuthorName  *string
	}

	err := pr.DB.WithContext(ctx).
		Select(`
			posts.id,
			posts.title,
			posts.slug,
			posts.content,
			posts.excerpt,
			posts.tags,
			posts.source,
			posts.access_level,
			posts.created_at,
			users.full_name as author_name
		`).
		Joins("left join users on users.id = posts.author_id").
		Where("id = ?", postID).
		Limit(1).
		Scan(&row).Error

	if err != nil {
		return nil, err
	}

	post := request.PostWithAuthor{
		ID:          row.ID,
		Title:       row.Title,
		Content:     row.Content,
		Excerpt:     row.Excerpt,
		CreatedAt:   row.CreatedAt,
		AccessLevel: row.AccessLevel,
		AuthorName:  row.AuthorName,
		TagsRaw:     row.TagsRaw,
	}

	
	if err := json.Unmarshal([]byte(row.TagsRaw), &post.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &post, nil
}