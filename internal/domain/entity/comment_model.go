package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

type Comment struct {
	ID        string   `gorm:"type:varchar(60);primaryKey"`
	PostID    string   `gorm:"type:varchar(60);not null"`
	Post      Post     `gorm:"foreignKey:PostID"`

	AuthorID  string   `gorm:"type:varchar(60);not null"`
	Author    User     `gorm:"foreignKey:AuthorID"`

	Content   string   `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Replies   []Reply  `gorm:"foreignKey:CommentID"`
}

/*
	Notes:

	Each comment belongs to a post and has an author.

	Comments can have many replies.
*/
