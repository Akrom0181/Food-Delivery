// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/Akrom0181/Food-Delivery/config"
	_ "github.com/Akrom0181/Food-Delivery/docs"
	"github.com/Akrom0181/Food-Delivery/internal/controller/http/v1/handler"
	"github.com/Akrom0181/Food-Delivery/internal/usecase"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

// NewRouter -.
// Swagger spec:
// @title       Food Delivery API
// @description This is a sample server Food Delivery server.
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis)

	// Initialize Casbin enforcer
	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e))

	// Swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1")

	user := v1.Group("/user")
	{
		user.POST("/", handlerV1.CreateUser)
		user.GET("/list", handlerV1.GetUsers)
		user.GET("/:id", handlerV1.GetUser)
		user.PUT("/", handlerV1.UpdateUser)
		user.DELETE("/:id", handlerV1.DeleteUser)
		user.POST("/upload", handlerV1.UploadProfilePic)
	}

	session := v1.Group("/session")
	{
		session.GET("/list", handlerV1.GetSessions)
		session.GET("/:id", handlerV1.GetSession)
		session.PUT("/", handlerV1.UpdateSession)
		session.DELETE("/:id", handlerV1.DeleteSession)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/logout", handlerV1.Logout)
		auth.POST("/register", handlerV1.Register)
		auth.POST("/verify-email", handlerV1.VerifyEmail)
		auth.POST("/login", handlerV1.Login)
	}

	report := v1.Group("/report")
	{
		report.POST("/", handlerV1.CreateReport)
		report.GET("/list", handlerV1.GetReports)
		report.GET("/:id", handlerV1.GetReport)
		report.PUT("/", handlerV1.UpdateReport)
		report.DELETE("/:id", handlerV1.DeleteReport)
	}

	notification := v1.Group("/notification")
	{
		notification.POST("/", handlerV1.CreateNotification)
		notification.GET("/list", handlerV1.GetNotifications)
		notification.GET("/:id", handlerV1.GetNotification)
		notification.PUT("/update-status", handlerV1.UpdateStatusNotification)
		notification.DELETE("/:id", handlerV1.DeleteNotification)
		notification.PUT("/:id", handlerV1.UpdateNotification)
	}

	firebase := v1.Group("/firebase")
	{
		firebase.POST("/", handlerV1.UploadFiles)
		firebase.DELETE("/:id", handlerV1.DeleteFile)
	}

	category := v1.Group("/category")
	{
		category.POST("/", handlerV1.CreateCategory)
		category.GET("/list", handlerV1.GetCategories)
		category.GET("/:id", handlerV1.GetCategory)
		category.PUT("/", handlerV1.UpdateCategory)
		category.DELETE("/:id", handlerV1.DeleteCategory)
	}

	product := v1.Group("/product")
	{
		product.POST("/", handlerV1.CreateProduct)
		product.GET("/list", handlerV1.GetProducts)
		product.GET("/:id", handlerV1.GetProduct)
		product.PUT("/", handlerV1.UpdateProduct)
		product.DELETE("/:id", handlerV1.DeleteProduct)
		product.PUT("/upload/:id", handlerV1.UploadProductPic)
	}

	banner := v1.Group("/banner")
	{
		banner.POST("/", handlerV1.UploadBanner)
		banner.GET("/list", handlerV1.GetBanners)
		banner.GET("/:id", handlerV1.GetBanner)
		banner.PUT("/", handlerV1.UpdateBanner)
		banner.DELETE("/:id", handlerV1.DeleteBanner)
	}
}
