package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

type Reply struct {
	ID         string   `gorm:"type:char(26);primaryKey"`
	CommentID  string   `gorm:"not null"`
	Comment    Comment     `gorm:"foreignKey:CommentID"`

	AuthorID   string   `gorm:"not null"`
	Author     User        `gorm:"foreignKey:AuthorID"`

	Content    string      `gorm:"type:text;not null"`

	CreatedAt  time.Time
	UpdatedAt  time.Time
}

/*
Notes:

Each reply belongs to a comment and has an author.

For nested replies later (replies to replies), you can add a ParentReplyID field.
*/