package server

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"pigpq/config"
	data "pigpq/database"
	"pigpq/internal/routers"
	"pigpq/internal/validator"
	"syscall"
	"time"
)

var (
	Cmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "go-layout server -c config.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			// 加载数据库配置
			data.InitDatabase()

			// 初始化验证器
			validator.InitValidatorTrans("zh")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	host string
	port int
)

func init() {
	Cmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "监听服务器地址")
	Cmd.Flags().IntVarP(&port, "port", "P", 0, "监听服务器端口")
}
func run() error {
	if port == 0 {
		port = config.Config.Port
	}

	// 初始化gin 路由
	engine := routers.SetRouters()

	// 创建 HTTP Server 实例
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: engine,
	}

	// 在 goroutine 中启动服务器
	go func() {
		zap.S().Infof("✅ API 服务器已启动，监听地址：%s:%d", host, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("❌ 服务启动失败: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.S().Info("🛑 收到退出信号，开始优雅关闭服务器...")

	// 创建 5 秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Infof("❌ 优雅关闭失败: %v", err)
		// 注意：即使 Shutdown 失败，我们也要继续退出
	} else {
		zap.S().Infof("👋 服务器已安全退出")
	}
	return nil
}
