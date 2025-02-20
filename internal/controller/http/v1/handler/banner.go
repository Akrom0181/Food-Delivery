package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/firebase"
	"github.com/gin-gonic/gin"
)

// GetBanner godoc
// @Router /banner/{id} [get]
// @Summary Get a banner by ID
// @Description Get a banner by ID
// @Security BearerAuth
// @Tags banner
// @Accept  json
// @Produce  json
// @Param id path string true "Banner ID"
// @Success 200 {object} entity.Banner
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBanner(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	banner, err := h.UseCase.BannerRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting banner") {
		return
	}

	ctx.JSON(200, banner)
}

// GetBanners godoc
// @Router /banner/list [get]
// @Summary Get a list of banners
// @Description Get a list of banners
// @Security BearerAuth
// @Tags banner
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BannerList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBanners(ctx *gin.Context) {
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
			Column: "title",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	banners, err := h.UseCase.BannerRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting banners") {
		return
	}

	ctx.JSON(200, banners)
}

// UpdateBanner godoc
// @Router /banner [put]
// @Summary Update a banner
// @Description Update a banner
// @Security BearerAuth
// @Tags banner
// @Accept  json
// @Produce  json
// @Param banner body entity.Banner true "Banner object"
// @Success 200 {object} entity.Banner
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBanner(ctx *gin.Context) {
	var (
		body entity.Banner
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	banner, err := h.UseCase.BannerRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating banner") {
		return
	}

	ctx.JSON(200, banner)
}

// DeleteBanner godoc
// @Router /banner/{id} [delete]
// @Summary Delete a banner
// @Description Delete a banner
// @Security BearerAuth
// @Tags banner
// @Accept  json
// @Produce  json
// @Param id path string true "Banner ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBanner(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") == "banner" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.BannerRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting banner") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Banner deleted successfully",
	})
}

// CreateBanner godoc
// @ID create_banner_pic_file
// @Router /banner [post]
// @Summary Upload a banner
// @Description Upload a banner
// @Security BearerAuth
// @Tags banner
// @Accept multipart/form-data
// @Produce json
// @Param file formData []file true "File to upload"
// @Param banner body entity.Banner true "Banner object"
// @Success 200 {object} entity.MultipleFileUploadResponse "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) UploadBanner(ctx *gin.Context) {
	var (
		body entity.Banner
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid file upload request", 400)
		return
	}

	if ctx.GetHeader("user_type") != "admin" {
		return
	}

	resp, err := firebase.UploadFiles(form)
	if h.HandleDbError(ctx, err, "Error uploading files") {
		return
	}

	body.ImageUrl = resp.Url[0].Url

	banner, err := h.UseCase.BannerRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating banner") {
		return
	}

	ctx.JSON(200, banner)
}
