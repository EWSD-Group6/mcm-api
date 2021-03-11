package user

import (
	"mcm-api/pkg/enforcer"
	"time"
)

type Entity struct {
	Id        int
	Name      string
	Email     string
	Password  string
	FacultyId *int
	Role      enforcer.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Entity) TableName() string {
	return "users"
}
