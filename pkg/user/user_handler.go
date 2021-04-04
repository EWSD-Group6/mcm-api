package user

import (
	"github.com/labstack/echo/v4"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/middleware"
	"net/http"
	"strconv"
)

type Handler struct {
	config  *config.Config
	service *Service
}

func NewUserHandler(config *config.Config, service *Service) *Handler {
	return &Handler{
		config:  config,
		service: service,
	}
}

func (h *Handler) Register(group *echo.Group) {
	group.Use(middleware.RequireAuthentication(h.config.JwtSecret))
	group.GET("", h.index, middleware.RequirePermission(enforcer.ReadUser))
	group.GET("/:id", h.getById, middleware.RequirePermission(enforcer.ReadUser))
	group.POST("", h.createUser, middleware.RequirePermission(enforcer.CreateUser))
	group.POST("/:id/status", h.updateStatus, middleware.RequirePermission(enforcer.UpdateUser))
	group.PATCH("/:id", h.updateUser, middleware.RequirePermission(enforcer.UpdateUser))
	group.DELETE("/:id", h.deleteUser, middleware.RequirePermission(enforcer.DeleteUser))
}

// @Tags Users
// @Summary List users
// @Description List users
// @Accept  json
// @Produce  json
// @Param params query user.UserIndexQuery false "user index query"
// @Success 200 {object} PaginateComposition
// @Security ApiKeyAuth
// @Router /users [get]
func (h *Handler) index(ctx echo.Context) error {
	query := new(UserIndexQuery)
	err := ctx.Bind(query)
	if err != nil {
		return err
	}
	loggedInUser, err := enforcer.GetLoggedInUser(ctx.Request().Context())
	if err != nil {
		return err
	}
	paginateRes, err := h.service.Find(ctx.Request().Context(), loggedInUser, query)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.JSON(http.StatusOK, paginateRes)
}

// @Tags Users
// @Summary Show a user
// @Description get user by ID
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} user.UserResponse
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *Handler) getById(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return apperror.
			New(apperror.ErrInvalid, "Id should be string", err).
			ToResponse(ctx)
	}
	res, err := h.service.FindById(ctx.Request().Context(), id)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.JSON(http.StatusOK, res)
}

// @Tags Users
// @Summary Create a user
// @Description Create a user
// @Accept  json
// @Produce  json
// @Param user body user.UserCreateReq true "create user"
// @Success 200 {object} user.UserResponse
// @Security ApiKeyAuth
// @Router /users [post]
func (h Handler) createUser(ctx echo.Context) error {
	req := &UserCreateReq{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	userResponse, err := h.service.CreateUser(ctx.Request().Context(), req)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.JSON(http.StatusCreated, userResponse)
}

// @Tags Users
// @Summary Update a user
// @Description Update a user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body user.UserUpdateReq true "update user"
// @Success 200 {object} user.UserResponse
// @Security ApiKeyAuth
// @Router /users/{id} [patch]
func (h Handler) updateUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return apperror.
			New(apperror.ErrInvalid, "Id should be string", err).
			ToResponse(ctx)
	}
	req := &UserUpdateReq{}
	err = ctx.Bind(req)
	if err != nil {
		return err
	}
	userResponse, err := h.service.Update(ctx.Request().Context(), id, req)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.JSON(http.StatusOK, userResponse)
}

// @Tags Users
// @Summary Delete a user
// @Description Delete a user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h Handler) deleteUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return apperror.
			New(apperror.ErrInvalid, "Id should be string", err).
			ToResponse(ctx)
	}
	err = h.service.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.NoContent(http.StatusOK)
}

// @Tags Users
// @Summary Update user status
// @Description Update user status
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Security ApiKeyAuth
// @Router /users/{id}/status [post]
func (h Handler) updateStatus(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return apperror.
			New(apperror.ErrInvalid, "Id should be string", err).
			ToResponse(ctx)
	}
	req := &UserUpdateStatus{}
	err = ctx.Bind(req)
	if err != nil {
		return err
	}
	err = h.service.ChangeStatus(ctx.Request().Context(), id, req)
	if err != nil {
		return apperror.HandleError(err, ctx)
	}
	return ctx.NoContent(http.StatusOK)
}
