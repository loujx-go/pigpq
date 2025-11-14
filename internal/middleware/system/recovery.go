package system

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"pigpq/config"
	e "pigpq/internal/pkg/errors"
	"pigpq/internal/pkg/response"
)

// CustomRecovery 自定义错误 (panic) 拦截中间件、对可能发生的错误进行拦截、统一记录
func CustomRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		errStr := ""
		if config.Config.Debug == true {
			errStr = fmt.Sprintf("%v", err)
		}
		response.NewResp().SetHttpCode(http.StatusInternalServerError).FailCode(c, e.ServerError, errStr)
	})
}

// PanicExceptionRecord  panic等异常记录
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	s1 := "An error occurred in the server's internal code："
	var build strings.Builder
	build.WriteString(s1)
	build.Write(b)
	errStr := build.String()
	zap.S().Error(errStr)
	return len(errStr), errors.New(errStr)
}
