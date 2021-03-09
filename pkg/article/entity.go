package article

import (
	"time"
)

type Entity struct {
	Id          int
	Title       string
	Description string
	Versions    []Version `gorm:"foreignKey:ArticleId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *Entity) TableName() string {
	return "articles"
}

type Version struct {
	Id           int
	Hash         string
	ArticleId    int
	LinkOriginal string
	LinkPdf      string
	CreatedAt    time.Time
}

func (v Version) TableName() string {
	return "article_versions"
}
