package comment

import (
	"github.com/labstack/echo/v4"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/middleware"
	"net/http"
)

type Handler struct {
	config  *config.Config
	service *Service
}

func NewHandler(config *config.Config, service *Service) *Handler {
	return &Handler{
		config:  config,
		service: service,
	}
}

func (h *Handler) Register(group *echo.Group) {
	group.Use(middleware.RequireAuthentication(h.config.JwtSecret))
	group.GET("", h.index, middleware.RequirePermission(enforcer.ReadComment))
	group.GET("/:id", h.getById, middleware.RequirePermission(enforcer.ReadComment))
	group.POST("", h.create, middleware.RequirePermission(enforcer.CreateComment))
	group.PUT("/:id", h.update, middleware.RequirePermission(enforcer.UpdateComment))
	group.DELETE("/:id", h.delete, middleware.RequirePermission(enforcer.DeleteComment))
}

// @Tags Comments
// @Summary List comments
// @Description List comments
// @Accept  json
// @Produce  json
// @Param params query comment.IndexQuery false "index query"
// @Success 200 {object} common.PaginateResponse{data=comment.CommentRes}
// @Security ApiKeyAuth
// @Router /comments [get]
func (h *Handler) index(context echo.Context) error {
	query := new(IndexQuery)
	err := context.Bind(query)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	paginateResponse, err := h.service.Find(context.Request().Context(), query)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, paginateResponse)
}

// @Tags Comments
// @Summary Show a comment
// @Description get comment by ID
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} comment.CommentRes
// @Security ApiKeyAuth
// @Router /comments/{id} [get]
func (h *Handler) getById(context echo.Context) error {
	result, err := h.service.FindById(context.Request().Context(), context.Param("id"))
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}

// @Tags Comments
// @Summary Create a comment
// @Description Create a comment
// @Accept  json
// @Produce  json
// @Param body body comment.CommentCreateReq true "create"
// @Success 200 {object} comment.CommentRes
// @Security ApiKeyAuth
// @Router /comments [post]
func (h *Handler) create(context echo.Context) error {
	body := new(CommentCreateReq)
	err := context.Bind(body)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	result, err := h.service.Create(context.Request().Context(), body)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}

// @Tags Comments
// @Summary Update a comment
// @Description Update a comment
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param body body comment.CommentUpdateReq true "update"
// @Success 200 {object} comment.CommentRes
// @Security ApiKeyAuth
// @Router /comments/{id} [put]
func (h *Handler) update(context echo.Context) error {
	body := new(CommentUpdateReq)
	err := context.Bind(body)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	result, err := h.service.Update(context.Request().Context(), context.Param("id"), body)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}

// @Tags Comments
// @Summary Delete a comment
// @Description Delete a comment
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200
// @Security ApiKeyAuth
// @Router /comments/{id} [delete]
func (h *Handler) delete(context echo.Context) error {
	err := h.service.Delete(context.Request().Context(), context.Param("id"))
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.NoContent(http.StatusNoContent)
}
