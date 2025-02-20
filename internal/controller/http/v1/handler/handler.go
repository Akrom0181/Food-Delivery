package handler

import (
	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/usecase"
	"github.com/Akrom0181/Food-Delivery/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
	}
}
