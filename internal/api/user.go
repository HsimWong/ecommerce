package api

import (
	"fmt"
	"time"

	"github.com/HsimWong/ecommerce/internal/service"
	"github.com/HsimWong/ecommerce/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserAPIHandler struct {
	userService *service.UserService
}

func NewUserAPIHandler() *UserAPIHandler {
	return &UserAPIHandler{
		userService: service.NewUserService(),
	}
}

func (uh *UserAPIHandler) Register(c *gin.Context) {
	logger.Log().Debug("Received Register request",
		zap.Time("timestamp", time.Now()))

	var req UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, err, "parameterError")
		return
	}

	logger.Log().Debug(fmt.Sprintf("%v", req))

}
