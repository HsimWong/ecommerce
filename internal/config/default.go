package config

import "go.uber.org/zap"

const (
	configFile      = "./configs/config.yaml"
	configType      = "yaml"
	LogFilePath     = "log/app.log"
	DefaultLogLevel = zap.DebugLevel
)
