package comment

import (
	"mcm-api/pkg/user"
	"time"
)

type Entity struct {
	Id             string `gorm:"default:null"`
	UserId         int
	User           user.Entity
	ContributionId int
	Content        string
	Resolved       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (e *Entity) TableName() string {
	return "comments"
}
