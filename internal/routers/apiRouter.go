package routers

import (
	"github.com/gin-gonic/gin"
	"pigpq/internal/controller"
	auth2 "pigpq/internal/controller/auth"
	"pigpq/internal/middleware"
)

func SetApiRoute(e *gin.Engine) {
	// version 1
	v1 := e.Group("api/v1")
	{
		demo := controller.NewDemoController()
		v1.GET("hello-world", demo.HelloWorld)

		// 登录
		auth := auth2.NewAuthController()
		v1.POST("login", auth.Login)

		/**============需要权限校验=============**/
		authReq := v1.Group("", middleware.Auth())
		{
			// 获取用户详情
			authReq.GET("user-info", auth.Info)
		}
	}
}
