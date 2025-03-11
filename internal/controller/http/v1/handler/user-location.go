package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateUserLocation godoc
// @Router /user/location [post]
// @Summary Create a new userlocation
// @Description Create a new userlocation
// @Security BearerAuth
// @Tags userlocation
// @Accept  json
// @Produce  json
// @Param userlocation body entity.UserLocation true "UserLocation object"
// @Success 201 {object} entity.UserLocation
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateUserLocation(ctx *gin.Context) {
	var (
		body entity.UserLocation
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.UserId = ctx.GetHeader("sub")

	userlocation, err := h.UseCase.UserLocationRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating user location") {
		return
	}

	ctx.JSON(201, userlocation)
}

// GetUserLocation godoc
// @Router /user/location/{id} [get]
// @Summary Get a userlocation by ID
// @Description Get a userlocation by ID
// @Security BearerAuth
// @Tags userlocation
// @Accept  json
// @Produce  json
// @Param id path string true "UserLocation ID"
// @Success 200 {object} entity.UserLocation
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetUserLocation(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	userlocation, err := h.UseCase.UserLocationRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting userlocation") {
		return
	}

	ctx.JSON(200, userlocation)
}

// GetUserLocations godoc
// @Router /user/location/list [get]
// @Summary Get a list of userlocations
// @Description Get a list of userlocations
// @Security BearerAuth
// @Tags userlocation
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param user_id query string false "user_id"
// @Success 200 {object} entity.ListUserLocation
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetUserLocations(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	userId := ctx.DefaultQuery("user_id", "")

	if ctx.GetHeader("user_type") == "user" {
		userId = ctx.GetHeader("sub")
	}

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "user_id",
			Type:   "eq",
			Value:  userId,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	userlocations, err := h.UseCase.UserLocationRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting user locations") {
		return
	}

	ctx.JSON(200, userlocations)
}

// UpdateUserLocation godoc
// @Router /user/location [put]
// @Summary Update a userlocation
// @Description Update a userlocation
// @Security BearerAuth
// @Tags userlocation
// @Accept  json
// @Produce  json
// @Param userlocation body entity.UserLocation true "UserLocation object"
// @Success 200 {object} entity.UserLocation
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateUserLocation(ctx *gin.Context) {
	var (
		body entity.UserLocation
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.UserId = ctx.GetHeader("sub")
	}

	userlocation, err := h.UseCase.UserLocationRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating userlocation") {
		return
	}

	ctx.JSON(200, userlocation)
}

// DeleteUserLocation godoc
// @Router /user/location/{id} [delete]
// @Summary Delete a userlocation
// @Description Delete a userlocation
// @Security BearerAuth
// @Tags userlocation
// @Accept  json
// @Produce  json
// @Param id path string true "UserLocation ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteUserLocation(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.UserLocationRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting user location") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "User location deleted successfully",
	})
}
