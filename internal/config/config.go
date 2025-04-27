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

type Database struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Dbname       string
	MaxOpenConns int
	MaxIdleConns int
	ConnTimeout  int
}

type Server struct {
	Addr string
	Port int
	Mode ServerMode
}

type Configuration struct {
	server   Server
	database Database
}

// Singleton config
var appConfig *Configuration
var once sync.Once

func initConfig(configfile, configType string) {
	viper.SetConfigFile(configfile) // 配置文件名 ()
	viper.SetConfigType(configType) // 配置文件类型

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
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

	// Important! using indirect assertion to avoid
	// pointer value change from outside
	tmpConfig := &struct {
		Server
		Database
	}{}
	if err := viper.Unmarshal(tmpConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		panic(err)
	}
	appConfig = &Configuration{
		server:   tmpConfig.Server,
		database: tmpConfig.Database,
	}
}

func assert(condition bool) { utils.Assert(condition) }

func (conf *Configuration) Validate() {
	assert(len(conf.server.Addr) > 0)
	assert(conf.server.Port <= 65535 && conf.server.Port > 0)
	assert(len(conf.database.Dbname) > 0)
	assert(len(conf.database.Host) > 0)
	assert(conf.database.Port <= 65535 && conf.database.Port > 0)
}

func (conf *Configuration) GetDBConfig() Database {
	return conf.database
}

func (conf *Configuration) GetServerConfig() Server {
	return conf.server
}

func Config(configPath ...string) *Configuration {
	once.Do(func() {
		initConfig(func() string {
			if len(configPath) <= 0 {
				return ConfigFile
			} else {
				if len(configPath[0]) <= 0 {
					return ConfigFile
				}
				return configPath[0]
			}
		}(), configType)
	})

	return appConfig
}
