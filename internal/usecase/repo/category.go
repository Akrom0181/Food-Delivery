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

type CategoryRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewCategoryRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *CategoryRepo {
	return &CategoryRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *CategoryRepo) Create(ctx context.Context, req entity.Category) (entity.Category, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("category").
		Columns(`id, name`).
		Values(req.Id, req.Name).ToSql()
	if err != nil {
		return entity.Category{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Category{}, err
	}

	return req, nil
}

func (r *CategoryRepo) GetSingle(ctx context.Context, req entity.CategorySingleRequest) (entity.Category, error) {
	response := entity.Category{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, created_at, updated_at`).
		From("category")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.Name != "":
		qeuryBuilder = qeuryBuilder.Where("name = ?", req.Name)
	default:
		return entity.Category{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Category{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.Name, &createdAt, &updatedAt)
	if err != nil {
		return entity.Category{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *CategoryRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.CategoryList, error) {
	var (
		response             = entity.CategoryList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, created_at, updated_at`).
		From("category")

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
		var item entity.Category
		err = rows.Scan(&item.Id, &item.Name, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("category").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *CategoryRepo) Update(ctx context.Context, req entity.Category) (entity.Category, error) {
	mp := map[string]interface{}{
		"name":       req.Name,
		"updated_at": "now()",
	}

	qeury, args, err := r.pg.Builder.Update("category").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.Category{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Category{}, err
	}

	return req, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("category").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("category").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
