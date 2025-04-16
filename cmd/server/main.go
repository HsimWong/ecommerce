package main

import (
	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/internal/router"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Config()
	appConfig.Validate()
	logger.Log().Debug("Server will be started started at ",
		zap.String("ServerAddr", appConfig.Server.Addr),
		zap.Int("ServerPort", appConfig.Server.Port),
		zap.String("dbhost", appConfig.Database.Host),
		zap.Int("dbport", appConfig.Database.Port),
	)

	r := router.NewRouter(config.SERVER_MODE_DEBUG)
	r.Run()
	select {}

}
