package app

import (
	"github.com/VmesteApp/auth-service/config"
	"github.com/VmesteApp/auth-service/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Info("logger init")
}
