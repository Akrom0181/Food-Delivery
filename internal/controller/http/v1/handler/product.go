package handler

import (
	"log"
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/firebase"
	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Router /product [post]
// @Summary Create a new product
// @Description Create a new product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param product body entity.Product true "Product object"
// @Success 201 {object} entity.Product
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var (
		body entity.Product
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	product, err := h.UseCase.ProductRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating product") {
		return
	}

	ctx.JSON(201, product)
}

// GetProduct godoc
// @Router /product/{id} [get]
// @Summary Get a product by ID
// @Description Get a product by ID
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetProduct(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	product, err := h.UseCase.ProductRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting product") {
		return
	}

	ctx.JSON(200, product)
}

// GetProducts godoc
// @Router /product/list [get]
// @Summary Get a list of products
// @Description Get a list of products
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.ProductList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetProducts(ctx *gin.Context) {
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

	products, err := h.UseCase.ProductRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting products") {
		return
	}

	ctx.JSON(200, products)
}

// UpdateProduct godoc
// @Router /product [put]
// @Summary Update a product
// @Description Update a product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param product body entity.Product true "Product object"
// @Success 200 {object} entity.Product
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var (
		body entity.Product
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") != "admin" {
		return
	}

	product, err := h.UseCase.ProductRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating product") {
		return
	}

	ctx.JSON(200, product)
}

// DeleteProduct godoc
// @Router /product/{id} [delete]
// @Summary Delete a product
// @Description Delete a product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("product_type") == "product" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.ProductRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting product") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// UploadProductPic godoc
// @ID upload_product_pic_file
// @Router /product/upload/{id} [put]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Security BearerAuth
// @Tags product
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param file formData []file true "File to upload"
// @Success 200 {object} entity.MultipleFileUploadResponse "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) UploadProductPic(ctx *gin.Context) {
	var (
		req = entity.Id{}
	)

	req.ID = ctx.Param("id")
	log.Println("req ----------------", req)

	form, err := ctx.MultipartForm()
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid file upload request", 400)
		return
	}

	if ctx.GetHeader("user_type") != "admin" {
		return
	}
	log.Println("product_id", req.ID)
	product, err := h.UseCase.ProductRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting product") {
		return
	}

	resp, err := firebase.UploadFiles(form)
	if h.HandleDbError(ctx, err, "Error uploading files") {
		return
	}
	log.Print("resp ------------------------", resp.Url[0].Url)
	_, err = h.UseCase.ProductRepo.Update(ctx, entity.Product{
		Id:          product.Id,
		CategoryId:  product.CategoryId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    resp.Url[0].Url,
	})

	if h.HandleDbError(ctx, err, "Error updating product") {
		return
	}

	ctx.JSON(200, resp)
}
