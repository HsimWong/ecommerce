package config

import "go.uber.org/zap"

const (
	configFile      = "./configs/config"
	configType      = "yaml"
	LogFilePath     = "log/app.log"
	DefaultLogLevel = zap.DebugLevel
)
