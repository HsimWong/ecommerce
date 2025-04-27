package database

import (
	"time"

	"gorm.io/gorm"
)

// Option 配置选项
type Option func(*DBManager)

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) Option {
	return func(m *DBManager) {
		m.maxRetries = maxRetries
	}
}

// WithRetryInterval 设置重试间隔
func WithRetryInterval(interval time.Duration) Option {
	return func(m *DBManager) {
		m.retryInterval = interval
	}
}

// WithGormConfig 设置GORM配置
func WithGormConfig(config *gorm.Config) Option {
	return func(m *DBManager) {
		m.gormConfig = config
	}
}
