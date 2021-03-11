package enforcer

import (
	"context"
	"github.com/labstack/echo/v4"
	"mcm-api/pkg/apperror"
)

type LoggedInUser struct {
	Id        int
	Email     string
	Name      string
	Role      Role
	FacultyId *int
}

const contextKey = "user"

func GetLoggedInUser(ctx context.Context) (*LoggedInUser, error) {
	if v, oke := ctx.Value(contextKey).(*LoggedInUser); oke {
		if oke {
			return v, nil
		}
	}
	return nil, apperror.New(apperror.ErrUnauthorized, "cant get logged in user", nil)
}

func SetLoggedInUser(ctx echo.Context, user *LoggedInUser) {
	ctx.SetRequest(ctx.Request().Clone(context.WithValue(ctx.Request().Context(), contextKey, user)))
}
