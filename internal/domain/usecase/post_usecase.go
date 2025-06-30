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
	"github.com/gosimple/slug"
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
	// Convert AuthorID if provided
	var authorID *ulid.ULID
	if req.AuthorID != nil {
		parsedID, err := ulid.Parse(*req.AuthorID)
		if err != nil {
			return err
		}
		authorID = &parsedID
	}

	// Set default access level if not provided
	if req.AccessLevel == nil {
		defaultLevel := "Subscribed"
		req.AccessLevel = &defaultLevel
	}

	// Marshal tags to JSON if provided
	var jsonTags *[]byte
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return err
		}
		jsonTags = &tagsBytes
	}

	// Making slug from title
	slug := slug.Make(req.Title)

	// Build post model
	post := &entity.Post{
		ID:          util.NewULID(),
		AuthorID:    authorID,
		Title:       req.Title,
		Slug:        slug,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Tags:        jsonTags,
		AccessLevel: *req.AccessLevel,
		Source:      req.Source,
	}

	// Create post in repository
	if err := usecase.Repo.Create(ctx, post); err != nil {
		return err
	}

	return nil
}


// Fetches posts
func (usecase *PostUsecase) FetchPosts(ctx context.Context, limit, offset int) ([]request.PostWithAuthor, int64, error) {
	return usecase.Repo.FetchPosts(ctx, limit, offset)
}


// Fetches a single post
func (usecase *PostUsecase) FetchPost(ctx context.Context, postID string) (*request.PostWithAuthor, error) {
	return usecase.Repo.FetchPost(ctx, postID)
}