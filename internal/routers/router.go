package routers

import (
	"io"
	"net/http"
	"pigpq/internal/pkg/errors"
	"pigpq/internal/pkg/response"

	"github.com/gin-gonic/gin"

	"pigpq/config"
	middleware "pigpq/internal/middleware/system"
)

func SetRouters() *gin.Engine {

	var engine *gin.Engine
	if config.Config.Debug == false {
		// 生产模式
		engine = ReleaseRouter()
		engine.Use(
			middleware.RequestCostHandler(),
			middleware.CustomLogger(),
			middleware.CustomRecovery(),
			middleware.CorsHandler(),
		)
	} else {
		// 开发调试模式
		engine = gin.Default()
		engine.Use(
			middleware.RequestCostHandler(),
			middleware.CustomRecovery(),
			middleware.CorsHandler(),
		)
	}

	// set up trusted agents
	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}
	// ping
	engine.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	SetApiRoute(engine)

	// 统一处理 404
	engine.NoRoute(func(c *gin.Context) {
		response.NewResp().SetHttpCode(http.StatusNotFound).FailCode(c, errors.NotFound)
	})

	return engine
}

// ReleaseRouter 生产模式使用官方建议设置为 release 模式
func ReleaseRouter() *gin.Engine {
	// 切换到生产模式
	gin.SetMode(gin.ReleaseMode)
	// 禁用 gin 输出接口访问日志
	gin.DefaultWriter = io.Discard

	engine := gin.New()

	return engine
}
