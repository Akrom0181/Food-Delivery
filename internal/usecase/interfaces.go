// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/Akrom0181/Food-Delivery/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.User) (entity.User, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// ReportRepo -.
	ReportRepoI interface {
		Create(ctx context.Context, req entity.Report) (entity.Report, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Report, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReportList, error)
		Update(ctx context.Context, req entity.Report) (entity.Report, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// NotificationRepo -.
	NotificationRepoI interface {
		Create(ctx context.Context, req entity.Notification) (entity.Notification, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Notification, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.NotificationList, error)
		Update(ctx context.Context, req entity.Notification) (entity.Notification, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateStatus(ctx context.Context, req entity.Notification) (entity.Notification, error)
	}

	// CategoryRepo -.
	CategoryRepoI interface {
		Create(ctx context.Context, req entity.Category) (entity.Category, error)
		GetSingle(ctx context.Context, req entity.CategorySingleRequest) (entity.Category, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.CategoryList, error)
		Update(ctx context.Context, req entity.Category) (entity.Category, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	ProductRepoI interface {
		Create(ctx context.Context, req entity.Product) (entity.Product, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Product, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ProductList, error)
		Update(ctx context.Context, req entity.Product) (entity.Product, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// BannerRepo -.
	BannerRepoI interface {
		Create(ctx context.Context, req entity.Banner) (entity.Banner, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Banner, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BannerList, error)
		Update(ctx context.Context, req entity.Banner) (entity.Banner, error)
		Delete(ctx context.Context, req entity.Id) error
	}
)
