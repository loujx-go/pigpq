package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pigpq/config"
	"pigpq/internal/model"
	"pigpq/internal/service/user"
	"strconv"
	"time"

	e "pigpq/internal/pkg/errors"
	"pigpq/internal/pkg/response"
	j "pigpq/internal/pkg/untils/jwt"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Fail(c, e.NotLogin, "未登录")
			c.Abort()
			return
		}

		accessToken, err := j.GetAccessToken(token)
		if err != nil {
			response.Fail(c, e.NotLogin, err.Error())
			c.Abort()
			return
		}

		// 解析Token
		customClaims := new(j.CustomClaims[model.User])
		result, err := j.Parse(accessToken, customClaims)
		if err != nil || result == nil {
			response.FailCode(c, e.NotLogin)
			c.Abort()
			return
		}

		uid, err := result.Claims.GetSubject()
		if err != nil {
			zap.L().Error("获取用户ID失败", zap.Error(err))
			response.FailCode(c, e.NotLogin)
			c.Abort()
			return
		}

		atoi, err := strconv.Atoi(uid)
		if err != nil {
			response.FailCode(c, e.NotLogin)
			c.Abort()
			return
		}

		exp, err := result.Claims.GetExpirationTime()
		// 获取过期时间返回err,或者exp为nil都返回错误
		if err != nil || exp == nil {
			response.FailCode(c, e.NotLogin)
			c.Abort()
			return
		}

		// 刷新时间大于0则判断剩余时间小于刷新时间时刷新Token并在Response header中返回
		if config.Config.Jwt.RefreshTTL > 0 {
			now := time.Now()
			diff := exp.Time.Sub(now)
			refreshTTL := config.Config.Jwt.RefreshTTL * time.Second
			fmt.Println(diff.Seconds(), refreshTTL)
			if diff < refreshTTL {
				tokenResponse, _ := user.NewUserService().Refresh(uint(atoi))
				c.Writer.Header().Set("refresh-access-token", tokenResponse.AccessToken)
				c.Writer.Header().Set("refresh-exp", strconv.FormatInt(tokenResponse.ExpiresAt, 10))
			}
		}

		c.Set("uid", uint(atoi))
		c.Next()
	}
}
