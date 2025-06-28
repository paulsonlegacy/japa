package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
	//"gorm.io/datatypes"
)

type Post struct {
	ID          ulid.ULID      `gorm:"type:char(26);primaryKey"`
	AuthorID    ulid.ULID      `gorm:"not null"`               // The user who authored the post
	Author      User           `gorm:"foreignKey:AuthorID"`    // Relation to User

	Title       string         `gorm:"not null"`
	Content     string         `gorm:"type:text;not null"`     // Full text content
	Excerpt     string         `gorm:"type:text"`              // Optional excerpt / preview
	Tags        []byte         `gorm:"type:json"`              // Array of tags (as JSON)
	IsPaidOnly  bool           `gorm:"default:false"`          // If true, only subscribed users can access

	CreatedAt   time.Time
	UpdatedAt   time.Time

	Comments    []Comment      `gorm:"foreignKey:PostID"`      // Relation to comments
}

/*
Tags used []byte but can use datatypes.JSON from "gorm.io/datatypes"
*/