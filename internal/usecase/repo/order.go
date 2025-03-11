package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	"github.com/Akrom0181/Food-Delivery/pkg/postgres"
	"github.com/google/uuid"
)

type OrderRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewOrderRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *OrderRepo {
	return &OrderRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *OrderRepo) Create(ctx context.Context, order entity.Order) (entity.Order, error) {

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return entity.Order{}, err
	}
	defer tx.Rollback(ctx)

	order.ID = uuid.NewString()
	orderQuery, orderArgs, err := r.pg.Builder.Insert("orders").
		Columns(`id, user_id, total_price, status, delivery_status, address, floor, door_number, entrance, latitude, longitude, branch_id`).
		Values(order.ID, order.UserID, order.TotalPrice, order.Status, order.DeliveryStatus, order.Address, order.Floor, order.DoorNumber, order.Entrance, order.Latitude, order.Longitude, order.BranchId).ToSql()
	if err != nil {
		return entity.Order{}, err
	}

	_, err = tx.Exec(ctx, orderQuery, orderArgs...)
	if err != nil {
		return entity.Order{}, err
	}

	for _, item := range order.OrderItems {
		item.Id = uuid.NewString()
		item.OrderId = order.ID
		itemQuery, itemArgs, err := r.pg.Builder.Insert("orderitems").
			Columns(`id, order_id, product_id, total_price, quantity, price`).
			Values(item.Id, item.OrderId, item.ProductId, item.TotalPrice, item.Quantity, item.Price).ToSql()
		if err != nil {
			return entity.Order{}, err
		}

		_, err = tx.Exec(ctx, itemQuery, itemArgs...)
		if err != nil {
			return entity.Order{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (r *OrderRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Order, error) {
	response := entity.Order{}
	var (
		createdAt, updatedAt time.Time
		courier_id           sql.NullString
	)

	// Query for the order details
	queryBuilder := r.pg.Builder.
		Select(`o.id, o.user_id, o.total_price, o.status, o.delivery_status, 
			o.address, o.floor, o.door_number, o.entrance, o.latitude, o.longitude, o.branch_id, o.courier_id, 
			o.created_at, o.updated_at`).
		From("orders AS o").
		Where("o.id = ?", req.ID)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Order{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&response.ID, &response.UserID, &response.TotalPrice, &response.Status, &response.DeliveryStatus,
		&response.Address, &response.Floor, &response.DoorNumber, &response.Entrance,
		&response.Latitude, &response.Longitude, &response.BranchId, &courier_id, &createdAt, &updatedAt,
	)
	if err != nil {
		return entity.Order{}, err
	}

	if courier_id.Valid {
		response.CourierId = courier_id.String
	} else {
		courier_id.String = ""
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Query for order items
	itemsQuery, itemsArgs, err := r.pg.Builder.
		Select(`id, order_id, product_id, total_price, quantity, price`).
		From("orderitems").
		Where("order_id = ?", req.ID).
		ToSql()
	if err != nil {
		return entity.Order{}, err
	}

	rows, err := r.pg.Pool.Query(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		return entity.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItems
		err := rows.Scan(&item.Id, &item.OrderId, &item.ProductId, &item.TotalPrice, &item.Quantity, &item.Price)
		if err != nil {
			return entity.Order{}, err
		}
		response.OrderItems = append(response.OrderItems, item)
	}

	return response, nil
}

func (r *OrderRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.OrderList, error) {
	var (
		response             = entity.OrderList{}
		createdAt, updatedAt time.Time
		courier_id           sql.NullString
	)

	queryBuilder := r.pg.Builder.
		Select(`o.id, o.user_id, o.total_price, o.status, o.delivery_status, o.address, o.floor, o.door_number, o.entrance, o.latitude, o.longitude, o.branch_id, o.courier_id, o.created_at, o.updated_at,
				oi.id, oi.order_id, oi.product_id, oi.total_price, oi.quantity, oi.price`).
		From("orders o").
		LeftJoin("orderitems oi ON o.id = oi.order_id")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var (
			orderItem entity.OrderItems
			order     entity.Order
		)
		err = rows.Scan(
			&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.DeliveryStatus,
			&order.Address, &order.Floor, &order.DoorNumber, &order.Entrance,
			&order.Latitude, &order.Longitude, &order.BranchId, &courier_id, &createdAt, &updatedAt,
			&orderItem.Id, &orderItem.OrderId, &orderItem.ProductId, &orderItem.TotalPrice,
			&orderItem.Quantity, &orderItem.Price,
		)
		if err != nil {
			return response, err
		}

		if courier_id.Valid {
			order.CourierId = courier_id.String
		} else {
			courier_id.String = ""
		}

		order.CreatedAt = createdAt.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)

		if existingOrder, exists := orderMap[order.ID]; exists {
			existingOrder.OrderItems = append(existingOrder.OrderItems, orderItem)
		} else {
			order.OrderItems = append(order.OrderItems, orderItem)
			orderMap[order.ID] = &order
		}
	}

	for _, order := range orderMap {
		response.Items = append(response.Items, *order)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(DISTINCT o.id)").From("orders o").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *OrderRepo) Update(ctx context.Context, req entity.Order) (entity.Order, error) {
	mp := map[string]interface{}{
		"id":              req.ID,
		"user_id":         req.UserID,
		"total_price":     req.TotalPrice,
		"status":          req.Status,
		"delivery_status": req.DeliveryStatus,
		"address":         req.Address,
		"floor":           req.Floor,
		"door_number":     req.DoorNumber,
		"entrance":        req.Entrance,
		"latitude":        req.Latitude,
		"longitude":       req.Longitude,
		"courier_id":      req.CourierId,
		"updated_at":      time.Now().In(time.UTC).Format(time.RFC3339),
	}

	query, args, err := r.pg.Builder.Update("orders").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Order{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Order{}, err
	}

	return req, nil
}

func (r *OrderRepo) Delete(ctx context.Context, req entity.Id) error {
	// Begin transaction
	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete order items
	orderItemsQuery, orderItemsArgs, err := r.pg.Builder.Delete("orderitems").Where("order_id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, orderItemsQuery, orderItemsArgs...)
	if err != nil {
		return err
	}

	// Delete order
	orderQuery, orderArgs, err := r.pg.Builder.Delete("orders").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, orderQuery, orderArgs...)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	query, args, err := r.pg.Builder.Update("orders").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}

func (r *OrderRepo) GetOrdersByBranch(ctx context.Context, req entity.GetListFilter) (entity.OrderList, error) {
	var (
		response             = entity.OrderList{}
		createdAt, updatedAt time.Time
		courier_id           sql.NullString
	)

	// Build base query
	queryBuilder := r.pg.Builder.
		Select(`o.id, o.user_id, o.total_price, o.status, o.delivery_status, o.address, 
				o.floor, o.door_number, o.entrance, o.latitude, o.longitude, o.branch_id, 
				o.courier_id, o.created_at, o.updated_at,
				oi.id, oi.order_id, oi.product_id, oi.total_price, oi.quantity, oi.price`).
		From("orders o").
		LeftJoin("orderitems oi ON o.id = oi.order_id")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	// Store orders in a map to group order items
	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var (
			orderItem entity.OrderItems
			order     entity.Order
		)
		err = rows.Scan(
			&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.DeliveryStatus,
			&order.Address, &order.Floor, &order.DoorNumber, &order.Entrance,
			&order.Latitude, &order.Longitude, &order.BranchId, &courier_id, &createdAt, &updatedAt,
			&orderItem.Id, &orderItem.OrderId, &orderItem.ProductId, &orderItem.TotalPrice,
			&orderItem.Quantity, &orderItem.Price,
		)
		if err != nil {
			return response, err
		}

		if courier_id.Valid {
			order.CourierId = courier_id.String
		} else {
			courier_id.String = ""
		}

		order.CreatedAt = createdAt.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)

		// Append order items to existing order in map
		if existingOrder, exists := orderMap[order.ID]; exists {
			existingOrder.OrderItems = append(existingOrder.OrderItems, orderItem)
		} else {
			order.OrderItems = append(order.OrderItems, orderItem)
			orderMap[order.ID] = &order
		}
	}

	// Convert map to list
	for _, order := range orderMap {
		response.Items = append(response.Items, *order)
	}

	// Get total count
	countQuery, args, err := r.pg.Builder.
		Select("COUNT(DISTINCT o.id)").
		From("orders o").
		Where(where).
		ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}
