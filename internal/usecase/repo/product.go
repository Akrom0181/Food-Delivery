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

type ProductRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewProductRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ProductRepo {
	return &ProductRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ProductRepo) Create(ctx context.Context, req entity.Product) (entity.Product, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("product").
		Columns(`id, category_id, name, description, price, image_url`).
		Values(req.Id, req.CategoryId, req.Name, req.Description, req.Price, req.ImageUrl).ToSql()
	if err != nil {
		return entity.Product{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Product{}, err
	}

	return req, nil
}

func (r *ProductRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Product, error) {
	response := entity.Product{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, category_id, name, description, price, image_url, created_at, updated_at`).
		From("product")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Product{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Product{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.CategoryId, &response.Name, &response.Description, &response.Price, &response.ImageUrl, &createdAt, &updatedAt)
	if err != nil {
		return entity.Product{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *ProductRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ProductList, error) {
	var (
		response             = entity.ProductList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, category_id, name, description, price, image_url, created_at, updated_at`).
		From("product")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Product
		err = rows.Scan(&item.Id, &item.CategoryId, &item.Name, &item.Description, &item.Price, &item.ImageUrl, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("product").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *ProductRepo) Update(ctx context.Context, req entity.Product) (entity.Product, error) {
	mp := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"image_url":   req.ImageUrl,
		"updated_at":  "now()",
	}

	qeury, args, err := r.pg.Builder.Update("product").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.Product{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Product{}, err
	}

	return req, nil
}

func (r *ProductRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("product").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("product").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
