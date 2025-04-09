package logger

import (
	"os"
	"sync"

	"github.com/HsimWong/ecommerce/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// type Logger struct {
// 	zapLogger *zap.Logger
// }

var ZapLogger *zap.Logger

// var LoggerInstance *Logger

var once sync.Once

func initLogger(path string) {
	writer := zapcore.AddSync(&lumberjack.Logger{
		// Filename:   "logs/app.log",
		Filename:   path,
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		config.DefaultLogLevel,
	)
	ZapLogger = zap.New(core, zap.AddCaller())
	// LoggerInstance = &Logger{
	// zapLogger: zap.New(core, zap.AddCaller()),
	// }
}

func Log() *zap.Logger {
	once.Do(func() { initLogger(config.LogFilePath) })
	return ZapLogger
}
