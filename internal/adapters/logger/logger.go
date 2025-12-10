package logger

import (
	"github.com/r6lik/recommendation_service/internal/adapters/config"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Server.Env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
