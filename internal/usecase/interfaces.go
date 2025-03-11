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

	// BranchRepo -.
	BranchRepoI interface {
		Create(ctx context.Context, req entity.Branch) (entity.Branch, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Branch, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BranchList, error)
		Update(ctx context.Context, req entity.Branch) (entity.Branch, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
		GetNearestBranch(ctx context.Context, lat, lon float64) (entity.Branch, error)
	}

	// UserLocationRepo -.
	UserLocationRepoI interface {
		Create(ctx context.Context, req entity.UserLocation) (entity.UserLocation, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.UserLocation, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ListUserLocation, error)
		Update(ctx context.Context, req entity.UserLocation) (entity.UserLocation, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	OrderRepoI interface {
		Create(ctx context.Context, req entity.Order) (entity.Order, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Order, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.OrderList, error)
		Update(ctx context.Context, req entity.Order) (entity.Order, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
		GetOrdersByBranch(ctx context.Context, req entity.GetListFilter) (entity.OrderList, error)
	}

	CourierRepoI interface {
		GetNearbyCouriers(ctx context.Context, lat, lng float64, radius float64) ([]entity.Courier, error)
		AssignOrderToCourier(ctx context.Context, orderID, courierID string) error
	}
)
