package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pigpq/internal/pkg/errors"
	"pigpq/internal/pkg/response"
)

type Api struct {
	errors.Error
}

// Success 业务成功响应
func (api Api) Success(c *gin.Context, data ...any) {
	r := response.NewResp()
	if data != nil {
		r.WithDataSuccess(c, data[0])
		return
	}
	r.Success(c)
}

// FailCode 业务失败响应
func (api Api) FailCode(c *gin.Context, code int, data ...any) {
	r := response.NewResp()
	if data != nil {
		r.WithData(data[0]).FailCode(c, code)
		return
	}
	r.FailCode(c, code)
}

// Fail 业务失败响应
func (api Api) Fail(c *gin.Context, code int, message string, data ...any) {
	r := response.NewResp()
	if data != nil {
		r.WithData(data[0]).FailCode(c, code, message)
		return
	}
	r.FailCode(c, code, message)
}

// Err 判断错误类型是自定义类型则自动返回错误中携带的code和message，否则返回服务器错误
func (api Api) Err(c *gin.Context, e error) {
	businessError, err := api.AsBusinessError(e)
	if err != nil {
		zap.S().Warn("Unknown error:", zap.Any("Error reason:", err))
		api.FailCode(c, errors.ServerError)
		return
	}

	api.Fail(c, businessError.GetCode(), businessError.GetMessage())
}
