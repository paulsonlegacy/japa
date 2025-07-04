package entity

import "time"


type ScrapedPost struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    Category    string    `gorm:"column:category;type:varchar(12);not null;default:'news'"` // news, learn, stories
    Title       string    `gorm:"column:title;type:text;not null"`
    ContentHTML string    `gorm:"column:content_html;type:longtext;not null"`     // longtext if using MySQL, text works too
    ContentText *string   `gorm:"column:content_text;type:longtext;null"`
    Excerpt     *string   `gorm:"column:excerpt;type:text;null"`
    PostImg     *string   `gorm:"column:post_img;type:text;null"`
    Source      string    `gorm:"column:source;type:text;uniqueIndex;not null"`
    Status      string    `gorm:"column:status;type:varchar(20);default:'pending';index"` // pending, published, ignored
    CreatedAt   time.Time `gorm:"autoCreateTime"` // PublishedAt
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
