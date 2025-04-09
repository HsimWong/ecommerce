package config

import "go.uber.org/zap"

const (
	configName      = "config"
	configType      = "yaml"
	configFilePath  = "./configs/"
	LogFilePath     = "log/app.log"
	DefaultLogLevel = zap.DebugLevel
)
