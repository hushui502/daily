package middleware

import (
	"rest-api/config"
	"rest-api/internal/auth"
	"rest-api/internal/session"
	"rest-api/pkg/logger"
)

type MiddlewareManager struct {
	sessUC session.UCSession
	authUC auth.UseCase
	cfg *config.Config
	origins []string
	logger logger.Logger
}

func NewMiddlewareManager(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
