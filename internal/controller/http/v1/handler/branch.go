package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Router /branch [post]
// @Summary Create a new branch
// @Description Create a new branch
// @Security BearerAuth
// @Tags branch
// @Accept  json
// @Produce  json
// @Param branch body entity.Branch true "Branch object"
// @Success 201 {object} entity.Branch
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBranch(ctx *gin.Context) {
	var (
		body entity.Branch
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	branch, err := h.UseCase.BranchRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating branch") {
		return
	}

	ctx.JSON(201, branch)
}

// GetBranch godoc
// @Router /branch/{id} [get]
// @Summary Get a branch by ID
// @Description Get a branch by ID
// @Security BearerAuth
// @Tags branch
// @Accept  json
// @Produce  json
// @Param id path string true "Branch ID"
// @Success 200 {object} entity.Branch
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBranch(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	branch, err := h.UseCase.BranchRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting branch") {
		return
	}

	ctx.JSON(200, branch)
}

// GetBranchs godoc
// @Router /branch/list [get]
// @Summary Get a list of branchs
// @Description Get a list of branchs
// @Security BearerAuth
// @Tags branch
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BranchList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBranches(ctx *gin.Context) {
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

	branchs, err := h.UseCase.BranchRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting branches") {
		return
	}

	ctx.JSON(200, branchs)
}

// UpdateBranch godoc
// @Router /branch [put]
// @Summary Update a branch
// @Description Update a branch
// @Security BearerAuth
// @Tags branch
// @Accept  json
// @Produce  json
// @Param branch body entity.Branch true "Branch object"
// @Success 200 {object} entity.Branch
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBranch(ctx *gin.Context) {
	var (
		body entity.Branch
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	branch, err := h.UseCase.BranchRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating branch") {
		return
	}

	ctx.JSON(200, branch)
}

// DeleteBranch godoc
// @Router /branch/{id} [delete]
// @Summary Delete a branch
// @Description Delete a branch
// @Security BearerAuth
// @Tags branch
// @Accept  json
// @Produce  json
// @Param id path string true "Branch ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBranch(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.BranchRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting branch") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Branch deleted successfully",
	})
}
