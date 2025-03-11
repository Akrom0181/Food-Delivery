package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	"github.com/Akrom0181/Food-Delivery/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type BranchRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBranchRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BranchRepo {
	return &BranchRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BranchRepo) Create(ctx context.Context, req entity.Branch) (entity.Branch, error) {
	req.Id = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("branch").
		Columns(`id, name, address, latitude, longitude, phone_number`).
		Values(req.Id, req.Name, req.Address, req.Latitude, req.Longitude, req.Phone).ToSql()
	if err != nil {
		return entity.Branch{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Branch{}, err
	}

	return req, nil
}

func (r *BranchRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Branch, error) {
	response := entity.Branch{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, name, address, latitude, longitude, phone_number, created_at, updated_at`).
		From("branch")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Branch{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Branch{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.Id, &response.Name, &response.Address, &response.Latitude, &response.Longitude, &response.Phone, &createdAt, &updatedAt)
	if err != nil {
		return entity.Branch{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BranchRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BranchList, error) {
	var (
		response             = entity.BranchList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, name, address, latitude, longitude, phone_number, created_at, updated_at`).
		From("branch")

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
		var item entity.Branch
		err = rows.Scan(&item.Id, &item.Name, &item.Address, &item.Latitude, &item.Longitude, &item.Phone, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("branch").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BranchRepo) Update(ctx context.Context, req entity.Branch) (entity.Branch, error) {
	mp := map[string]interface{}{
		"name":         req.Name,
		"address":      req.Address,
		"latitude":     req.Latitude,
		"longitude":    req.Longitude,
		"phone_number": req.Phone,
		"updated_at":   "now()",
	}

	query, args, err := r.pg.Builder.Update("branch").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.Branch{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Branch{}, err
	}

	return req, nil
}

func (r *BranchRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("branch").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BranchRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	query, args, err := r.pg.Builder.Update("branch").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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

func (r *BranchRepo) GetNearestBranch(ctx context.Context, lat, lon float64) (entity.Branch, error) {
	const searchRadius = 10000.0 // 10 km radius ichidagi filiallarni qidiramiz
	var (
		created_at, updated_at time.Time
	)
	var branch entity.Branch
	query := `
		SELECT id, name, address, latitude, longitude, phone_number, created_at, updated_at
		FROM branch
		WHERE earth_distance(ll_to_earth($1, $2), ll_to_earth(latitude, longitude)) < $3
		ORDER BY earth_distance(ll_to_earth($1, $2), ll_to_earth(latitude, longitude))
		LIMIT 1
	`

	err := r.pg.Pool.QueryRow(ctx, query, lat, lon, searchRadius).Scan(
		&branch.Id, &branch.Name, &branch.Address, &branch.Latitude, &branch.Longitude,
		&branch.Phone, &created_at, &updated_at,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Branch{}, errors.New("there is no branch near you in 10 km radius")
		}
		return entity.Branch{}, err
	}

	branch.CreatedAt = created_at.Format(time.RFC3339)
	branch.UpdatedAt = updated_at.Format(time.RFC3339)
	return branch, nil
}
