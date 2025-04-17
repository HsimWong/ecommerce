package main

import (
	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/internal/database"
	"github.com/HsimWong/ecommerce/internal/router"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Config()
	appConfig.Validate()
	logger.Log().Debug("Server will be started started at ",
		zap.String("ServerAddr", appConfig.GetServerConfig().Addr),
		zap.Int("ServerPort", appConfig.GetServerConfig().Port),
		zap.String("dbhost", appConfig.GetDBConfig().Host),
		zap.Int("dbport", appConfig.GetDBConfig().Port),
	)

	database.DBConn(appConfig.GetDBConfig())

	r := router.NewRouter(config.SERVER_MODE_DEBUG)
	go r.Run()

	select {}

}
