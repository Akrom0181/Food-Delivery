package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Router /order [post]
// @Summary Create a new order
// @Description Create a new order
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param order body entity.Order true "Order object"
// @Success 201 {object} entity.Order
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateOrder(ctx *gin.Context) {
	var body entity.Order

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	branch, err := h.UseCase.BranchRepo.GetNearestBranch(ctx, body.Latitude, body.Longitude)
	if h.HandleDbError(ctx, err, "Error getting nearest branch") {
		return
	}

	body.BranchId = branch.Id
	body.UserID = ctx.GetHeader("sub")

	// Calculate total price
	var totalPrice float64
	for i, item := range body.OrderItems {
		if item.Quantity == 0 {
			h.ReturnError(ctx, config.ErrorBadRequest, "Quantity must be greater than 0", 400)
			return
		}
		product, err := h.UseCase.ProductRepo.GetSingle(ctx, entity.Id{ID: item.ProductId})
		if h.HandleDbError(ctx, err, "Error getting product") {
			return
		}
		body.OrderItems[i].Price = product.Price
		body.OrderItems[i].TotalPrice = product.Price * float64(item.Quantity)
		totalPrice += body.OrderItems[i].TotalPrice
	}
	body.TotalPrice = totalPrice
	body.Status = "pending"

	order, err := h.UseCase.OrderRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating order") {
		return
	}

	ctx.JSON(201, order)
}

// GetOrder godoc
// @Router /order/{id} [get]
// @Summary Get a order by ID
// @Description Get a order by ID
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetOrder(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	order, err := h.UseCase.OrderRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting order") {
		return
	}

	ctx.JSON(200, order)
}

// GetOrders godoc
// @Router /order/list [get]
// @Summary Get a list of orders
// @Description Get a list of orders
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param branch_id query string true "branch_id"
// @Success 200 {object} entity.OrderList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetOrders(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("branch_id", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "branch_id",
			Type:   "eq",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	orders, err := h.UseCase.OrderRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting orders") {
		return
	}

	ctx.JSON(200, orders)
}

// UpdateOrder godoc
// @Router /order [put]
// @Summary Update a order
// @Description Update a order
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param order body entity.Order true "Order object"
// @Success 200 {object} entity.Order
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateOrder(ctx *gin.Context) {
	var (
		body entity.Order
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.UserID = ctx.GetHeader("sub")
	}

	getorder, err := h.UseCase.OrderRepo.GetSingle(ctx, entity.Id{ID: body.ID})
	if h.HandleDbError(ctx, err, "Error getting order") {
		return
	}

	if body.Status == "cancelled" && getorder.Status == "picked_up" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Order already picked up and cannot be cancelled", 400)
		return
	}

	order, err := h.UseCase.OrderRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating order") {
		return
	}

	ctx.JSON(200, order)
}

// DeleteOrder godoc
// @Router /order/{id} [delete]
// @Summary Delete a order
// @Description Delete a order
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteOrder(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") != "admin" {
		return
	}

	err := h.UseCase.OrderRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting order") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Order deleted successfully",
	})
}

// GetBranchOrders godoc
// @Router /order/bybranch [get]
// @Summary Get a list of branch orders
// @Description Get a list of branch orders
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param page query number true "Page number"
// @Param limit query number true "Number of results per page"
// @Param branch_id query string true "Branch ID"
// @Param status query string true "Order status"
// @Param delivery_status query string true "Delivery status"
// @Success 200 {object} entity.OrderList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBranchOrders(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	branchId := ctx.DefaultQuery("branch_id", "")
	status := ctx.DefaultQuery("status", "")
	delivery_s := ctx.DefaultQuery("delivery_status", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)

	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "branch_id",
			Type:   "eq",
			Value:  branchId,
		},
		entity.Filter{
			Column: "status",
			Type:   "eq",
			Value:  status,
		},
		entity.Filter{
			Column: "delivery_status",
			Type:   "eq",
			Value:  delivery_s,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	orders, err := h.UseCase.OrderRepo.GetOrdersByBranch(ctx, req)
	if err != nil {
		h.HandleDbError(ctx, err, "Error getting orders")
		return
	}

	ctx.JSON(200, orders)
}
