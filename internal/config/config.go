package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/HsimWong/ecommerce/internal/utils"
	"github.com/spf13/viper"
)

type ServerMode string

const (
	SERVER_MODE_RELEASE ServerMode = "release"
	SERVER_MODE_DEBUG   ServerMode = "debug"
	SERVER_MODE_TEST    ServerMode = "test"
)

type Configuration struct {
	Server struct {
		Addr string
		Port int
		Mode ServerMode
	}
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Dbname   string
	}
}

// Singleton config
var AppConfig *Configuration
var once sync.Once

func initConfig(configfile, configType string) {
	viper.SetConfigFile(configfile) // 配置文件名 ()
	viper.SetConfigType(configType) // 配置文件类型

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}
	AppConfig = &Configuration{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}
	// 特别处理密码，优先从环境变量获取
	if viper.GetString("database.postgres.password") == "${DB_PASSWORD}" {
		if pwd := os.Getenv("DB_PASSWORD"); pwd != "" {
			viper.Set("database.postgres.password", pwd)
		} else {
			panic(fmt.Errorf("DB_PASSWORD environment variable is required"))
		}
	}
}

func assert(condition bool) { utils.Assert(condition) }

func (conf *Configuration) Validate() {
	assert(len(conf.Server.Addr) > 0)
	assert(conf.Server.Port <= 65535)
	assert(len(conf.Database.Dbname) > 0)
	assert(len(conf.Database.Host) > 0)
	assert(conf.Database.Port < 65535)
}

func Config(configPath ...string) *Configuration {
	once.Do(func() {
		initConfig(func() string {
			if len(configPath) <= 0 {
				return configFile
			} else {
				return configPath[0]
			}
		}(), configType)
	})

	return AppConfig
}
