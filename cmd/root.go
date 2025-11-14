package cmd

import (
	"fmt"
	"os"
	"pigpq/cmd/cron"
	"pigpq/cmd/server"
	"pigpq/config"
	"pigpq/pkg/logger"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "go-pig",
		Short: "pig",
		Long:  `gin`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 初始化配置
			config.InitConfig(configPath)
			// 时区配置
			if config.Config.Timezone != nil {
				location, err := time.LoadLocation(*config.Config.Timezone)
				if err != nil {
					fmt.Println("Error loading location:", err)
					return
				}
				time.Local = location
			}
			// 加载logger
			logger.InitLogger()
		},
		// 如果有相关的 action 要执行，请取消下面这行代码的注释
		// Run: func(cmd *cobra.Command, args []string) { },
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有传入子命令，默认执行 server.Cmd
			if len(args) == 0 {
				// 将当前命令的 flag 传递给 server.Cmd
				server.Cmd.SetArgs(args)
				// 继承 PersistentFlags（如 --config）
				server.Cmd.PersistentFlags().AddFlagSet(cmd.PersistentFlags())

				if err := server.Cmd.Execute(); err != nil {
					os.Exit(1)
				}
				return
			}

			// 否则显示帮助
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config/config.yaml", "The absolute path of the configuration file")
	// 启动服务 go-layout server
	rootCmd.AddCommand(cron.Cmd)

}

// Execute 将所有子命令添加到root命令并适当设置标志。
// 这由 main.main() 调用。它只需要对 rootCmd 调用一次。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}

	// ✅ 程序退出前刷新日志
	if logger.Logger != nil {
		_ = logger.Logger.Sync()
	}
}
