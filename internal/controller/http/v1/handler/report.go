package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateReport godoc
// @Router /report [post]
// @Summary Create a new report
// @Description Create a new report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param report body entity.Report true "Report object"
// @Success 201 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateReport(ctx *gin.Context) {
	var (
		body entity.Report
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.UserID = ctx.GetHeader("sub")

	report, err := h.UseCase.ReportRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating report") {
		return
	}

	ctx.JSON(201, report)
}

// GetUser godoc
// @Router /report/{id} [get]
// @Summary Get a report by ID
// @Description Get a report by ID
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param id path string true "Report ID"
// @Success 200 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReport(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	report, err := h.UseCase.ReportRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting report") {
		return
	}

	ctx.JSON(200, report)
}

// GetUsers godoc
// @Router /report/list [get]
// @Summary Get a list of reports
// @Description Get a list of reports
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.ReportList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReports(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "user_id",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "business_id",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "reason",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.ReportRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting reports") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateReport godoc
// @Router /report [put]
// @Summary Update a report
// @Description Update a report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param report body entity.Report true "Report object"
// @Success 200 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateReport(ctx *gin.Context) {
	var (
		body entity.Report
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.ID = ctx.GetHeader("sub")
	}

	report, err := h.UseCase.ReportRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating report") {
		return
	}

	ctx.JSON(200, report)
}

// DeleteUser godoc
// @Router /report/{id} [delete]
// @Summary Delete a report
// @Description Delete a report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param id path string true "Report ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteReport(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") == "user" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.ReportRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting user") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Report deleted successfully",
	})
}
