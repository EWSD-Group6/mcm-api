package user

import (
	"mcm-api/pkg/enforcer"
	"time"
)

type UserStatus string

const (
	UserActive  UserStatus = "active"
	UserDisable UserStatus = "disable"
)

type Entity struct {
	Id        int
	Name      string
	Email     string
	Password  string
	FacultyId *int
	Role      enforcer.Role
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Entity) TableName() string {
	return "users"
}
