package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Router /category [post]
// @Summary Create a new category
// @Description Create a new category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param category body entity.Category true "Category object"
// @Success 201 {object} entity.Category
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateCategory(ctx *gin.Context) {
	var (
		body entity.Category
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") != "admin" {
		return
	}

	category, err := h.UseCase.CategoryRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating category") {
		return
	}

	ctx.JSON(201, category)
}

// GetCategory godoc
// @Router /category/{id} [get]
// @Summary Get a category by ID
// @Description Get a category by ID
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.Category
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetCategory(ctx *gin.Context) {
	var (
		req entity.CategorySingleRequest
	)

	req.ID = ctx.Param("id")

	category, err := h.UseCase.CategoryRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting category") {
		return
	}

	ctx.JSON(200, category)
}

// GetCategorys godoc
// @Router /category/list [get]
// @Summary Get a list of categorys
// @Description Get a list of categorys
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.CategoryList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetCategories(ctx *gin.Context) {
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
			Column: "name",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	categorys, err := h.UseCase.CategoryRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting categorys") {
		return
	}

	ctx.JSON(200, categorys)
}

// UpdateCategory godoc
// @Router /category [put]
// @Summary Update a category
// @Description Update a category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param category body entity.Category true "Category object"
// @Success 200 {object} entity.Category
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var (
		body entity.Category
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "admin" {
		body.Id = ctx.GetHeader("sub")
	}

	category, err := h.UseCase.CategoryRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating category") {
		return
	}

	ctx.JSON(200, category)
}

// DeleteCategory godoc
// @Router /category/{id} [delete]
// @Summary Delete a category
// @Description Delete a category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") != "admin" {
		return
	}

	err := h.UseCase.CategoryRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting category") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Category deleted successfully",
	})
}
