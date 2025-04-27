package main

import (
	"flag"
	"fmt"

	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/internal/database"
	"github.com/HsimWong/ecommerce/internal/init/dbinit"
	"github.com/HsimWong/ecommerce/internal/router"

	"github.com/HsimWong/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

var (
	configPath *string
	configMode *bool
)

func handleFlags() {
	configPath = flag.String(
		"config",
		config.ConfigFile,
		fmt.Sprintf("Config file path, default: %s", config.ConfigFile),
	)
	configMode = flag.Bool(
		"init",
		false,
		"Config mode flag, specifies wither to config initial running, default false",
	)

	flag.Parse()
	fmt.Printf("configMode: %v", configMode)
}

func main() {
	handleFlags()

	fmt.Print("=================================================")
	fmt.Printf("configMode: %v", configMode)
	appConfig := config.Config(*configPath)
	appConfig.Validate()
	// db := database.DBConn(appConfig.GetDBConfig())
	db, err := database.GetDatabaseInstance()
	if err != nil {
		logger.Log().Error(err.Error())
	}
	if *configMode {
		logger.Log().Info("Running in configure mode")
		gm, err := db.GetDB()
		if err != nil {
			logger.Log().Fatal("Gorm fetching failed")
		}
		err = dbinit.CreateTables(gm)
		if err != nil {
			logger.Log().Fatal("Database init failed")
		}
		logger.Log().Info("Database init success")
		return
	} else {
		logger.Log().Debug("Server will be started started at ",
			zap.String("ServerAddr", appConfig.GetServerConfig().Addr),
			zap.Int("ServerPort", appConfig.GetServerConfig().Port),
			zap.String("dbhost", appConfig.GetDBConfig().Host),
			zap.Int("dbport", appConfig.GetDBConfig().Port),
		)

		r := router.NewRouter(config.SERVER_MODE_DEBUG)
		go r.Run()

		select {}
	}

}
