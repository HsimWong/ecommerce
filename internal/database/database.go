package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *DBManager
	once     sync.Once
	mu       sync.Mutex
)

// DBManager 数据库管理器(单例)
type DBManager struct {
	cfg             config.Database
	gormConfig      *gorm.Config
	mu              sync.RWMutex
	db              *gorm.DB
	reconnecting    bool
	reconnectNotify chan struct{}
	closeCh         chan struct{}
	maxRetries      int
	retryInterval   time.Duration
}

// GetInstance 获取DBManager单例实例

func GetDatabaseInstance(opts ...Option) (*DBManager, error) {
	if instance != nil {
		return instance, nil
	}

	mu.Lock()
	defer mu.Unlock()

	// 再次检查，防止在获取锁期间其他goroutine已经创建了实例
	if instance != nil {
		return instance, nil
	}

	cfg := config.Config().GetDBConfig()

	// 使用sync.Once确保初始化只执行一次
	var initErr error
	once.Do(func() {
		instance = &DBManager{
			cfg:             cfg,
			gormConfig:      &gorm.Config{},
			reconnectNotify: make(chan struct{}),
			closeCh:         make(chan struct{}),
			maxRetries:      100,             // 默认重试次数
			retryInterval:   5 * time.Second, // 默认重试间隔
		}

		// 应用选项
		for _, opt := range opts {
			opt(instance)
		}

		// 初始连接
		if err := instance.connect(); err != nil {
			initErr = fmt.Errorf("initial connection failed: %w", err)
			instance = nil // 初始化失败，重置instance
			return
		}

		// 启动健康检查协程
		go instance.healthCheck()
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

// buildDSN 构建PostgreSQL连接字符串
func (m *DBManager) buildDSN() string {
	cfg := m.cfg
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname)
}

// connect 建立数据库连接
func (m *DBManager) connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var db *gorm.DB
	var err error

	// 带超时的连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.cfg.ConnTimeout)*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		db, err = gorm.Open(postgres.Open(m.buildDSN()), m.gormConfig)
	}()

	select {
	case <-done:
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}

		// 配置连接池
		if m.cfg.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(m.cfg.MaxOpenConns)
		}
		if m.cfg.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(m.cfg.MaxIdleConns)
		}
		sqlDB.SetConnMaxLifetime(time.Hour)

		m.db = db
		return nil
	case <-ctx.Done():
		return fmt.Errorf("connection timeout after %d seconds", m.cfg.ConnTimeout)
	}
}

// GetDB 获取数据库连接
func (m *DBManager) GetDB() (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	// 检查连接是否有效
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		// 触发异步重连
		go m.reconnect()
		return nil, fmt.Errorf("database connection is down: %w", err)
	}

	return m.db, nil
}

// reconnect 重新连接数据库
func (m *DBManager) reconnect() {
	m.mu.Lock()
	if m.reconnecting {
		m.mu.Unlock()
		return
	}
	m.reconnecting = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.reconnecting = false
		m.mu.Unlock()
	}()

	for i := 0; i < m.maxRetries; i++ {
		select {
		case <-m.closeCh:
			return // 如果管理器已关闭，则退出
		default:
		}

		if err := m.connect(); err == nil {
			m.mu.Lock()
			notify := m.reconnectNotify
			m.reconnectNotify = make(chan struct{})
			m.mu.Unlock()

			// 通知所有等待的调用者
			close(notify)
			return
		}

		time.Sleep(m.retryInterval)
	}
}

// healthCheck 定期健康检查
func (m *DBManager) healthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.closeCh:
			return
		case <-ticker.C:
			m.mu.RLock()
			db := m.db
			m.mu.RUnlock()

			if db == nil {
				continue
			}

			sqlDB, err := db.DB()
			if err != nil {
				continue
			}

			if err := sqlDB.Ping(); err != nil {
				go m.reconnect()
			}
		}
	}
}

// Close 关闭数据库连接
func (m *DBManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-m.closeCh:
		return nil // 已经关闭
	default:
		close(m.closeCh)
	}

	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// WaitForReconnect 等待重连完成
func (m *DBManager) WaitForReconnect(ctx context.Context) error {
	m.mu.RLock()
	notify := m.reconnectNotify
	m.mu.RUnlock()

	select {
	case <-notify:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
