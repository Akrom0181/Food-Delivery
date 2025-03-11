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

type BannerRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBannerRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BannerRepo {
	return &BannerRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BannerRepo) Create(ctx context.Context, req entity.Banner) (entity.Banner, error) {
	req.Id = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("banner").
		Columns(`id, title, image_url`).
		Values(req.Id, req.Title, req.ImageUrl).ToSql()
	if err != nil {
		return entity.Banner{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Banner{}, err
	}

	return req, nil
}

func (r *BannerRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Banner, error) {
	response := entity.Banner{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, banner, created_at, updated_at`).
		From("banner")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Banner{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Banner{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.Id, &response.ImageUrl, &createdAt, &updatedAt)
	if err != nil {
		return entity.Banner{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BannerRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BannerList, error) {
	var (
		response             = entity.BannerList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, title, image_url, created_at, updated_at`).
		From("banner")

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
		var item entity.Banner
		err = rows.Scan(&item.Id, &item.ImageUrl, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("banner").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BannerRepo) Update(ctx context.Context, req entity.Banner) (entity.Banner, error) {
	mp := map[string]interface{}{
		"title":      req.Title,
		"image_url":  req.ImageUrl,
		"updated_at": "now()",
	}

	query, args, err := r.pg.Builder.Update("banner").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.Banner{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Banner{}, err
	}

	return req, nil
}

func (r *BannerRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("banner").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
