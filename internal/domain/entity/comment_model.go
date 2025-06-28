package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Comment struct {
	ID        ulid.ULID   `gorm:"type:char(26);primaryKey"`
	PostID    ulid.ULID   `gorm:"not null"`
	Post      Post        `gorm:"foreignKey:PostID"`

	AuthorID  ulid.ULID   `gorm:"not null"`
	Author    User        `gorm:"foreignKey:AuthorID"`

	Content   string      `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Replies   []Reply     `gorm:"foreignKey:CommentID"`
}

/*
Notes:

Each comment belongs to a post and has an author.

Comments can have many replies.
*/
