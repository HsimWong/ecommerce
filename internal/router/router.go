package router

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/HsimWong/ecommerce/internal/api"
	"github.com/HsimWong/ecommerce/internal/config"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	r        *gin.Engine
	apiGroup *gin.RouterGroup
}

func (rt *Router) init() {
	rt.r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	apiGroup := rt.r.Group("/api/v1")

	// Grouping codes logically
	{
		userGroup := apiGroup.Group("/user")
		userHandler := api.NewUserAPIHandler()
		userGroup.POST("/register", userHandler.Register)
	}

	rt.apiGroup = apiGroup
}

func NewRouter(serverMode config.ServerMode) *Router {
	logger.Log()
	gin.SetMode(map[config.ServerMode]string{
		config.SERVER_MODE_RELEASE: gin.ReleaseMode,
		config.SERVER_MODE_DEBUG:   gin.DebugMode,
		config.SERVER_MODE_TEST:    gin.TestMode,
	}[serverMode])

	r := gin.New()
	r.SetTrustedProxies([]string{"192.168.0.0/16", "10.0.0.0/8"})

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

	rt := &Router{
		r: r,
	}
	rt.init()
	return rt
}

func (r *Router) Run() {
	cfg := config.Config()
	start := time.Now()
	if cfg == nil {
		panic("Running un-configged router")
	}

	logger.Log().Info("Trying to start HTTP Server",
		zap.String("ServerAddr", cfg.GetServerConfig().Addr),
		zap.Int("ServerPort", cfg.GetServerConfig().Port))

	if err := r.r.Run(fmt.Sprintf("%s:%d", cfg.GetServerConfig().Addr, cfg.GetServerConfig().Port)); err != nil {
		logger.Log().Fatal("HTTP Server failed",
			zap.Time("EndTime", time.Now()),
			zap.Duration("Duration", time.Since(start)),
		)
	}

}
