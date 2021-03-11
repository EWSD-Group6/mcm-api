package user

import (
	"mcm-api/pkg/common"
	"mcm-api/pkg/enforcer"
)

type UserResponse struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	FacultyId *int          `json:"facultyId"`
	Role      enforcer.Role `json:"role"`
	common.TrackTime
}
