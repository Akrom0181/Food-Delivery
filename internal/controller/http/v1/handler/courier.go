package handler

// import (
// 	"log"

// 	"github.com/Akrom0181/Food-Delivery/config"
// 	"github.com/Akrom0181/Food-Delivery/internal/entity"
// 	"github.com/gin-gonic/gin"
// )

// // AssignCourier godoc
// // @Router /assign-courier/{order_id} [post]
// // @Summary Assign a courier to an order
// // @Description Retrieves an order by ID, finds the nearest available courier, and assigns the order
// // @Security BearerAuth
// // @Tags courier
// // @Accept  json
// // @Produce  json
// // @Param order_id path string true "Order ID"
// // @Success 200 {object} map[string]string
// // @Failure 400 {object} entity.ErrorResponse
// // @Failure 403 {object} entity.ErrorResponse
// // @Failure 404 {object} entity.ErrorResponse
// // @Failure 500 {object} entity.ErrorResponse
// func (h *Handler) AssignCourierToOrder(ctx *gin.Context) {
// 	orderID := ctx.Param("order_id") // Get order ID from URL path
// 	id := entity.Id{ID: orderID}
// 	// Ensure only admins can assign couriers
// 	if ctx.GetHeader("user_type") != "admin" {
// 		h.ReturnError(ctx, config.ErrorForbidden, "Permission denied", 403)
// 		return
// 	}

// 	// Get order details
// 	order, err := h.UseCase.OrderRepo.GetSingle(ctx, id)
// 	if err != nil {
// 		h.ReturnError(ctx, config.ErrorNotFound, "Order not found", 404)
// 		return
// 	}

// 	// Define the search radius (e.g., 5000 meters = 5km)
// 	const searchRadius = 5000.0

// 	// Find nearby couriers
// 	couriers, err := h.UseCase.CourierRepo.GetNearbyCouriers(ctx, order.Latitude, order.Longitude, searchRadius)
// 	if err != nil || len(couriers) == 0 {
// 		h.ReturnError(ctx, config.ErrorNotFound, "No available couriers nearby", 404)
// 		return
// 	}

// 	// Pick the first available courier
// 	courier := couriers[0]

// 	// Assign order
// 	err = h.UseCase.CourierRepo.AssignOrderToCourier(ctx, order.ID, courier.ID)
// 	if h.HandleDbError(ctx, err, "Error assigning courier") {
// 		return
// 	}

// 	// TODO: Send push notification to courier here

// 	log.Printf("Order %s assigned to courier %s", order.ID, courier.ID)

// 	ctx.JSON(200, gin.H{
// 		"message":    "Courier assigned successfully",
// 		"courier_id": courier.ID,
// 	})
// }
