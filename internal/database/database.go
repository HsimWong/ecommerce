package database

import (
	"fmt"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// type Database interface {
// 	*sqlx.DB
// 	Run()
// }

type Database struct {
	*sqlx.DB
}

// 全局变量（小写字母开头，限制包内访问）
var dbinstance Database

// InitPostgres 初始化PostgreSQL连接池
func InitPostgres(cfg config.Database) (Database, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname,
	)
	var err error
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 连接池配置
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// func (db *sqlx.DB) Run() {

// }

// GetDB 获取已初始化的连接池（避免重复初始化）
func GetDB() *sqlx.DB {
	if db == nil {
		panic("database not initialized. Call InitPostgres() first")
	}
	return db
}
