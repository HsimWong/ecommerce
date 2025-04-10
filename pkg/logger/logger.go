package logger

import (
	"os"
	"sync"

	"github.com/HsimWong/ecommerce/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ZapLogger *zap.Logger
var loggerMutex *sync.RWMutex

var loggerDefaultConfig = &lumberjack.Logger{
	MaxSize:    100, // MB
	MaxBackups: 30,
	MaxAge:     30, // days
	Compress:   true,
}

var once sync.Once

func initLogger(path string) error {
	loggerMutex = &sync.RWMutex{}

	lock := loggerMutex.TryLock()
	if !lock {
		return ErrSetLogLevelFailed
	}
	defer loggerMutex.Unlock()
	loggerDefaultConfig.Filename = path
	writer := zapcore.AddSync(loggerDefaultConfig)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		config.DefaultLogLevel,
	)
	ZapLogger = zap.New(core, zap.AddCaller())
	return nil
}

func SetLevel(level zapcore.Level) error {
	// With write lock contained, the func can be called in Log() and SetLevel()
	once.Do(func() { initLogger(config.LogFilePath) })
	lock := loggerMutex.TryLock()
	if !lock {
		return ErrSetLogLevelFailed
	}
	defer loggerMutex.Unlock()

	writer := zapcore.AddSync(loggerDefaultConfig)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		level,
	)

	ZapLogger = zap.New(core, zap.AddCaller())
	return nil
}

func Log() *zap.Logger {
	once.Do(func() { initLogger(config.LogFilePath) })
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return ZapLogger
}
