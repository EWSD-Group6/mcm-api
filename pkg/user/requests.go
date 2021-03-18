package user

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"mcm-api/pkg/common"
	"mcm-api/pkg/enforcer"
)

type UserIndexQuery struct {
	Role enforcer.Role `query:"role" enums:"admin,marketing_manager,marketing_coordinator,student,guest"`
	common.PaginateQuery
}

type UserCreateReq struct {
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	Role      enforcer.Role `json:"role" enums:"admin,marketing_manager,marketing_coordinator,student,guest"`
	FacultyId *int          `json:"facultyId"`
}

func (c *UserCreateReq) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&c.Role, validation.Required, validation.In(
			enforcer.Guest, enforcer.Student, enforcer.MarketingCoordinator, enforcer.MarketingManager, enforcer.Administrator)),
		validation.Field(&c.FacultyId, validation.Required.When(isRoleRequiredFaculty(c.Role))),
	)
}

type UserUpdateReq struct {
	Name      *string        `json:"name"`
	Email     *string        `json:"email"`
	Password  *string        `json:"password"`
	Role      *enforcer.Role `json:"role"`
	FacultyId *int           `json:"facultyId"`
}

func (c UserUpdateReq) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Length(5, 50)),
		validation.Field(&c.Email, is.Email),
		validation.Field(&c.Password, validation.Length(5, 50)),
		validation.Field(&c.Role, validation.In(
			enforcer.Guest, enforcer.Student, enforcer.MarketingCoordinator, enforcer.MarketingManager, enforcer.Administrator)),
		validation.Field(&c.FacultyId),
	)
}

func isRoleRequiredFaculty(role enforcer.Role) bool {
	switch role {
	case enforcer.Administrator:
		return false
	case enforcer.MarketingManager:
		return false
	case enforcer.MarketingCoordinator:
		return true
	case enforcer.Student:
		return true
	case enforcer.Guest:
		return true
	default:
		return true
	}
}
