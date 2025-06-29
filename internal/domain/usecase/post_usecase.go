package usecase

import (
	"context"
	//"errors"
	"encoding/json"
	//"fmt"
	//"time"

	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	//"japa/internal/infrastructure/mail"
	//"japa/internal/pkg"
	"japa/internal/util"

	//"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/oklog/ulid/v2"
)

// TYPES

// PostUsecase handles user-related business logic
type PostUsecase struct {
	Repo      *repository.PostRepository
	DB        *gorm.DB
}

// Initialize PostUsecase
func NewPostUsecase(repo *repository.PostRepository, db *gorm.DB) *PostUsecase {
	return &PostUsecase{Repo: repo, DB: db}
}

// Creates new post
func (usecase *PostUsecase) CreatePost(ctx context.Context, req request.CreatePostRequest) error {
	// Type conversions
	authorID, err := ulid.Parse(req.AuthorID)
	if err != nil {
		return err
	}

	jsonTags, err := json.Marshal(req.Tags)
	if err != nil {
		return err
	}

	// Model build 
	post := &entity.Post{
		ID:          util.NewULID(),
		AuthorID:    authorID,
		Title:       req.Title,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Tags:        jsonTags,
		AccessLevel: *req.AccessLevel,
	}

	// create post via repository
	if err := usecase.Repo.Create(ctx, post); err != nil {
		return err 
	}

	return nil
}

// Fetches posts
func (usecase *PostUsecase) FetchPosts(ctx context.Context, limit, offset int) ([]request.PostWithAuthor, int64, error) {
	return usecase.Repo.Fetch(ctx, limit, offset)
}