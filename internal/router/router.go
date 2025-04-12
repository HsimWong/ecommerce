package router

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	r        *gin.Engine
	apiGroup *gin.RouterGroup
}

func NewRouter(serverMode config.ServerMode) *Router {
	gin.SetMode(map[config.ServerMode]string{
		config.SERVER_MODE_RELEASE: gin.ReleaseMode,
		config.SERVER_MODE_DEBUG:   gin.DebugMode,
		config.SERVER_MODE_TEST:    gin.TestMode,
	}[serverMode])
	r := gin.New()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		logger.Log().Error("HTTP Panic",
			zap.Any("error", err),
			zap.String("path", c.Request.URL.Path),
			zap.ByteString("stack", debug.Stack()))

		c.AbortWithStatus(500)

	}))

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.Log().Debug("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()))
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	apiGroup := r.Group("/api/v1")

	return &Router{
		r:        r,
		apiGroup: apiGroup,
	}
}

func (r *Router) Run() {
	cfg := config.Config()
	start := time.Now()
	if cfg == nil {
		panic("Running un-configged router")
	}

	logger.Log().Info("Trying to start HTTP Server",
		zap.String("ServerAddr", cfg.Server.Addr),
		zap.Int("ServerPort", cfg.Database.Port))

	if err := r.r.Run(fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)); err != nil {
		logger.Log().Fatal("HTTP Server failed",
			zap.String("EndTime", time.Now().String()),
			zap.String("Duration", time.Since(start).String()),
		)
	}
}
