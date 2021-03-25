package article

import (
	"time"
)

type Entity struct {
	Id        int
	Versions  []*Version `gorm:"foreignKey:ArticleId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Entity) TableName() string {
	return "articles"
}

type Version struct {
	Id           int       `json:"id"`
	Hash         string    `json:"hash"`
	ArticleId    int       `json:"articleId"`
	LinkOriginal string    `json:"linkOriginal"`
	LinkPdf      string    `json:"linkPdf"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (v Version) TableName() string {
	return "article_versions"
}
