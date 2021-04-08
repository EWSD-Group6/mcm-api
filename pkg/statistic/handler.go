package statistic

import (
	"github.com/labstack/echo/v4"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
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
	group.GET("/admin-dashboard", h.adminDashboard)
	group.GET("/contribution-faculty-chart", h.contributionFacultyChart)
	group.GET("/contribution-student-chart", h.contributionStudentChart)
}

// @Tags Statistics
// @Summary Admin Dashboard Data
// @Description Admin Dashboard Data
// @Accept  json
// @Produce  json
// @Success 200 {object} statistic.AdminDashboard
// @Security ApiKeyAuth
// @Router /statistics/admin-dashboard [get]
func (h *Handler) adminDashboard(context echo.Context) error {
	result, err := h.service.adminDashboard(context.Request().Context())
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}

// @Tags Statistics
// @Summary Contribution group by faculty data
// @Description Contribution group by faculty data
// @Accept  json
// @Produce  json
// @Param params query statistic.ContributionFacultyChartQuery false "query"
// @Success 200 {object} statistic.ContributionFacultyChart
// @Security ApiKeyAuth
// @Router /statistics/contribution-faculty-chart [get]
func (h Handler) contributionFacultyChart(context echo.Context) error {
	query := new(ContributionFacultyChartQuery)
	err := context.Bind(query)
	if err != nil {
		return err
	}
	result, err := h.service.contributionFacultyChart(context.Request().Context(), query)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}

// @Tags Statistics
// @Summary Contribution group by student data
// @Description Contribution group by student data
// @Accept  json
// @Produce  json
// @Param params query statistic.ContributionStudentChartQuery false "query"
// @Success 200 {object} statistic.ContributionStudentChart
// @Security ApiKeyAuth
// @Router /statistics/contribution-student-chart [get]
func (h Handler) contributionStudentChart(context echo.Context) error {
	query := new(ContributionStudentChartQuery)
	err := context.Bind(query)
	if err != nil {
		return err
	}
	result, err := h.service.contributionStudentChart(context.Request().Context(), query)
	if err != nil {
		return apperror.HandleError(err, context)
	}
	return context.JSON(http.StatusOK, result)
}
