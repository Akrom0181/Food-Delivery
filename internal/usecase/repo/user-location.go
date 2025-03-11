package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	"github.com/Akrom0181/Food-Delivery/pkg/postgres"
	"github.com/google/uuid"
)

type UserLocationRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserLocationRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserLocationRepo {
	return &UserLocationRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserLocationRepo) Create(ctx context.Context, req entity.UserLocation) (entity.UserLocation, error) {
	req.Id = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("user_location").
		Columns(`id, user_id, address, entrance, floor, door_number, latitude, longitude`).
		Values(req.Id, req.UserId, req.Address, req.Entrance, req.Floor, req.DoorNumber, req.Latitude, req.Longitude).ToSql()
	if err != nil {
		return entity.UserLocation{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.UserLocation{}, err
	}

	return req, nil
}

func (r *UserLocationRepo) GetSingle(ctx context.Context, req entity.Id) (entity.UserLocation, error) {
	response := entity.UserLocation{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, address, entrance, floor, door_number, latitude, longitude, created_at, updated_at`).
		From("user_location")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	default:
		return entity.UserLocation{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.UserLocation{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.Id, &response.UserId, &response.Address, &response.Entrance, &response.Floor, &response.DoorNumber, &response.Latitude, &response.Longitude, &createdAt, &updatedAt)
	if err != nil {
		return entity.UserLocation{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *UserLocationRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ListUserLocation, error) {
	var (
		response             = entity.ListUserLocation{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, address, entrance, floor, door_number, latitude, longitude, created_at, updated_at`).
		From("user_location")

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

	for rows.Next() {
		var item entity.UserLocation
		err = rows.Scan(&item.Id, &item.UserId, &item.Address, &item.Entrance, &item.Floor, &item.DoorNumber, &item.Latitude, &item.Longitude, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("user_location").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UserLocationRepo) Update(ctx context.Context, req entity.UserLocation) (entity.UserLocation, error) {
	mp := map[string]interface{}{
		"user_id":     req.UserId,
		"address":     req.Address,
		"entrance":    req.Entrance,
		"floor":       req.Floor,
		"door_number": req.DoorNumber,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
		"updated_at":  "now()",
	}

	query, args, err := r.pg.Builder.Update("user_location").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.UserLocation{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.UserLocation{}, err
	}

	return req, nil
}

func (r *UserLocationRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("user_location").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
