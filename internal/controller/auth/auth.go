package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pigpq/config"
	. "pigpq/internal/controller"
	"pigpq/internal/model"
	"pigpq/internal/pkg/response"
	"pigpq/internal/request/user"
	userservice "pigpq/internal/service/user"
	"pigpq/internal/validator"
)

type AuthController struct {
	Api
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (api AuthController) Login(c *gin.Context) {
	// 验证请求参数
	var login user.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		msg := validator.TranslateError(err)
		response.NewResp().Fail(c, 400, msg)
		return
	}

	// 开发环境不验证验证码
	if config.Config.AppEnv != "local" {
		// 验证验证码
	}

	// 用户登录
	result, err := userservice.NewUserService().Login(login.Phone)
	if err != nil {
		response.NewResp().Fail(c, 400, err.Error())
		return
	}

	api.Success(c, result)
}

// 获取用户详情
func (api AuthController) Info(c *gin.Context) {
	userIdAny, exits := c.Get("uid")
	if !exits {
		api.Fail(c, 400, "用户不存在")
		return
	}

	zap.S().Infof("useridany:%v", userIdAny.(uint))
	userId, ok := userIdAny.(uint)
	if !ok {
		api.Fail(c, 400, "用户ID格式错误")
		return
	}

	result := model.NewUser().GetUserById(userId)

	api.Success(c, result)
}
