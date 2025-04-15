package api

import (
	"github.com/HsimWong/ecommerce/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ResponseError(c *gin.Context, err error, msg string) {
	logger.Log().Error(msg, zap.Error(err))
	c.JSON(400, gin.H{"error": msg + " " + err.Error()})
}
