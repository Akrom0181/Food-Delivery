package repo

import (
	"context"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	"github.com/Akrom0181/Food-Delivery/pkg/postgres"
)

type CourierRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewCourierRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *CourierRepo {
	return &CourierRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *CourierRepo) GetNearbyCouriers(ctx context.Context, lat, lng float64, radius float64) ([]entity.Courier, error) {
	query := `
		SELECT id, user_id, status, latitude, longitude, last_updated 
		FROM couriers 
		WHERE status = 'active' 
		AND ST_DistanceSphere(
			ST_MakePoint(longitude, latitude),
			ST_MakePoint($1, $2)
		) <= $3
		ORDER BY ST_DistanceSphere(
			ST_MakePoint(longitude, latitude),
			ST_MakePoint($1, $2)
		) ASC
		LIMIT 1;
	`
	rows, err := r.pg.Pool.Query(ctx, query, lng, lat, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []entity.Courier
	for rows.Next() {
		var courier entity.Courier
		if err := rows.Scan(&courier.ID, &courier.UserID, &courier.Status, &courier.Latitude, &courier.Longitude, &courier.LastUpdated); err != nil {
			return nil, err
		}
		couriers = append(couriers, courier)
	}

	return couriers, nil
}

func (r *CourierRepo) AssignOrderToCourier(ctx context.Context, orderID, courierID string) error {
	query := `
		INSERT INTO order_courier (order_id, courier_id, assigned_at, status)
		VALUES ($1, $2, NOW(), 'pending')
	`
	_, err := r.pg.Pool.Exec(ctx, query, orderID, courierID)
	if err != nil {
		return err
	}
	return nil
}
