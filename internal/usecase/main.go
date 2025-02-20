package usecase

import (
	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/usecase/repo"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	"github.com/Akrom0181/Food-Delivery/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo         UserRepoI
	SessionRepo      SessionRepoI
	ReportRepo       ReportRepoI
	NotificationRepo NotificationRepoI
	CategoryRepo     CategoryRepoI
	ProductRepo      ProductRepoI
	BannerRepo       BannerRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:         repo.NewUserRepo(pg, config, logger),
		SessionRepo:      repo.NewSessionRepo(pg, config, logger),
		ReportRepo:       repo.NewReportRepo(pg, config, logger),
		NotificationRepo: repo.NewNotificationRepo(pg, config, logger),
		CategoryRepo:     repo.NewCategoryRepo(pg, config, logger),
		ProductRepo:      repo.NewProductRepo(pg, config, logger),
		BannerRepo:       repo.NewBannerRepo(pg, config, logger),
	}
}
