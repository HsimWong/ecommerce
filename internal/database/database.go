package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Database struct {
	*sqlx.DB
	cfg      config.Database
	poollock sync.RWMutex
}

// 全局变量（小写字母开头，限制包内访问）
var dbinstance *Database
var once sync.Once

// InitPostgres 初始化PostgreSQL连接池
func initPostgres(cfg config.Database) error {
	defer func() {
		dbinstance.Run()
	}()
	logger.Log().Info("initting database connection",
		zap.String("host", cfg.Host),
	)

	if dbinstance == nil {
		dbinstance = &Database{poollock: sync.RWMutex{}, cfg: cfg}
		logger.Log().Info("config loaded")
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname,
	)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(cfg.ConnTimeout))
	defer cancel()

	dbchan := make(chan *sqlx.DB, 1)
	var err error
	go func() {
		var db *sqlx.DB
		db, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %w", err)
		}
		dbchan <- db

	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("database connection timeout")
	case db := <-dbchan:
		if db != nil {
			// 连接池配置
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetConnMaxLifetime(5 * time.Minute)

			dbinstance.Load(db)

			return nil
		} else {
			return err
		}

	}
}

func (db *Database) Run() {

	go func(db *Database) {
		defer func(cfg config.Database) {
			if db.DB != nil {
				db.Close()
			}
			initPostgres(cfg)
		}(db.cfg)
		var err error
		for {
			if db.DB == nil {
				logger.Log().Error("database connection not established")
				// time.Sleep(5 * time.Second)
				return
			}
			// retry for 5 times in each loop
			for i := 0; i < 5; i++ {
				err = db.Ping()
				if err == nil {
					break
				}
				logger.Log().Debug("trying reconnecting database",
					zap.Int("trial", i),
					zap.Error(err),
				)
				time.Sleep(1 * time.Second)
			}

			if err != nil {
				logger.Log().Info("database connection offline,",
					zap.Error(err),
				)
				return
			}
			time.Sleep(5 * time.Second)
		}
	}(db)
}

func (db *Database) Load(dbload *sqlx.DB) {
	db.poollock.Lock()
	defer db.poollock.Unlock()
	db.DB = dbload
}

// GetDB 获取已初始化的连接池（避免重复初始化）
func DBConn(cfg ...config.Database) *Database {
	once.Do(func() {
		initPostgres(cfg[0])
	})
	return dbinstance
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	db.poollock.RLock()
	defer db.poollock.RUnlock()
	if db.DB != nil {
		return db.DB.Exec(query, args...)
	}
	return nil, fmt.Errorf("database lost connection")

}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	db.poollock.RLock()
	defer db.poollock.RUnlock()
	if db.DB != nil {
		return db.DB.Query(query, args...)
	}
	return nil, fmt.Errorf("database lost connection")
}
