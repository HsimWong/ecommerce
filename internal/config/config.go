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

func initConfig(configfile, configType string) {
	viper.SetConfigFile(configfile) // 配置文件名 ()
	viper.SetConfigType(configType) // 配置文件类型

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
		initConfig(configFile, configType)
	})
	return &AppConfig
}
