package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Addr string
		Port int
	}
	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		Dbname   string
	}
}

// Singleton config
var AppConfig Config
var once sync.Once

func initConfig(configName, configType, filepath string) {
	viper.SetConfigName(configName) // 配置文件名 (不带扩展名)
	viper.SetConfigType(configType) // 配置文件类型
	viper.AddConfigPath(filepath)   // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}
}

func Configure() *Config {
	once.Do(func() {
		initConfig(configName, configType, configFilePath)
	})
	return &AppConfig
}
