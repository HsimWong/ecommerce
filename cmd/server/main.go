package main

import (
	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Configure()
	logger.Log().Debug("Server will be started started at ",
		zap.String("ServerAddr", appConfig.Server.Addr),
		zap.Int("ServerPort", appConfig.Server.Port),
		zap.String("dbhost", appConfig.Database.Host),
		zap.Int("dbport", appConfig.Database.Port),
	)
}
