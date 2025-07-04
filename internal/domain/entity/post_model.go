package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
	//"gorm.io/datatypes"
)

type Post struct {
	ID          ulid.ULID      `gorm:"type:char(26);primaryKey"`
	AuthorID    *ulid.ULID      `gorm:"column:author_id;null"`  // The user who authored the post
	Author      *User           `gorm:"foreignKey:AuthorID"`    // Relation to User

	Category    string         `gorm:"column:category;type:varchar(12);not null;default:'learn'"` // news, stories, learn
	Title       string         `gorm:"column:title;type:text;not null"`
	Slug        string         `gorm:"column:slug;type:text;uniqueIndex;not null"`
	Content     string         `gorm:"column:content;type:text;not null"`      // Full text content
	Excerpt     *string        `gorm:"column:excerpt;type:text;null"`        // Optional excerpt / preview
	Tags        *[]byte        `gorm:"column:tags;type:json;null"`        // Array of tags (as JSON)

	// Source
	Source      *string        `gorm:"column:source;type:text;null"`

	// Audience controlSubscribed
	AccessLevel  string         `gorm:"column:access_level;not null;default:Subscribed"` // Public, Registered & Subscribed

	CreatedAt    time.Time      `gorm:"autoCreateTime"`  // PublishedAt
	UpdatedAt    time.Time	    `gorm:"autoUpdateTime"`

	Comments     []Comment      `gorm:"foreignKey:PostID"`      // Relation to comments
}

/*
Tags used []byte but can use datatypes.JSON from "gorm.io/datatypes"
*/