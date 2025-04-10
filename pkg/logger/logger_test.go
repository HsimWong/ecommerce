package logger

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(m *testing.T) {
	config.Configure("../../configs/config.yaml")

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 30; i++ {
			time.Sleep(1 * time.Second)
			Log().Debug("Debugging",
				zap.Int("CurrentLoop", i),
			)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(5 * time.Second)
			SetLevel([2]zapcore.Level{zap.InfoLevel, zap.DebugLevel}[i%2])
			fmt.Printf("set level to %v\n", [2]zapcore.Level{zap.InfoLevel, zap.DebugLevel}[i%2])
		}
		wg.Done()
	}()
	wg.Wait()

}
