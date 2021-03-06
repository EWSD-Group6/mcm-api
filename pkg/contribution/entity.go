package contribution

import (
	"mcm-api/pkg/article"
	"mcm-api/pkg/user"
	"time"
)

type Status string

const (
	Accepted  Status = "accepted"
	Rejected  Status = "rejected"
	Reviewing Status = "reviewing"
)

type Entity struct {
	Id                  int
	UserId              int
	User                user.Entity `gorm:"foreignKey:UserId"`
	ContributeSessionId int
	ArticleId           *int
	Article             *article.Entity `gorm:"foreignKey:ArticleId"`
	Title               string
	Description         string
	Status              Status
	Images              []ImageEntity `gorm:"foreignKey:ContributionId"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (e Entity) TableName() string {
	return "contributions"
}

type ImageEntity struct {
	Key            string `gorm:"primaryKey"`
	ContributionId int
	Title          string
}

func (i ImageEntity) TableName() string {
	return "images"
}
