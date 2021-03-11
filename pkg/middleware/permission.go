package middleware

import (
	"github.com/labstack/echo/v4"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/enforcer"
)

func RequirePermission(permission enforcer.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			user, err := enforcer.GetLoggedInUser(context.Request().Context())
			if err != nil {
				return apperror.HandleError(err, context)
			}
			if !enforcer.Can(user.Role, permission) {
				return apperror.HandleError(
					apperror.New(
						apperror.ErrForbidden,
						"insufficient permission",
						nil),
					context,
				)
			}
			return next(context)
		}
	}
}
