package main

import (
	"fmt"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/internal/database"
	"github.com/HsimWong/ecommerce/internal/router"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Config()
	// appConfig.GetServerConfig().Port = 6060
	appConfig.Validate()
	logger.Log().Debug("Server will be started started at ",
		zap.String("ServerAddr", appConfig.GetServerConfig().Addr),
		zap.Int("ServerPort", appConfig.GetServerConfig().Port),
		zap.String("dbhost", appConfig.GetDBConfig().Host),
		zap.Int("dbport", appConfig.GetDBConfig().Port),
	)

	r := router.NewRouter(config.SERVER_MODE_DEBUG)
	go r.Run()

	db, err := database.InitPostgres(appConfig.GetDBConfig())
	if err != nil {
		logger.Log().Fatal("Database init failed", zap.Error(err))
	}
	for {
		logger.Log().Info(fmt.Sprintf("%v", db.Ping()))
		time.Sleep(5 * time.Second)
	}

	select {}

}
