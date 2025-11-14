package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

var (
	Cmd = &cobra.Command{
		Use:     "cron",
		Short:   "Starting a scheduled task",
		Example: "go-layout cron",
		PreRun: func(cmd *cobra.Command, args []string) {
			// 计划任务中使用数据请先初始化数据库链接
			//data.InitData()
		},
		Run: func(cmd *cobra.Command, args []string) {
			Start()
		},
	}
)

func Start() {
	myLog := myLogger{}
	crontab := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(myLog)))
	job := cron.NewChain(cron.SkipIfStillRunning(myLog), cron.Recover(myLog)).Then(cron.FuncJob(func() {
		zap.S().Infof("%s:%s\n", time.Now().Format("2006-01-02 15:04:05"), " 计划任务 demo 执行")
	}))
	_, err := crontab.AddJob("0/5 * * * * *", job)

	if err != nil {
		panic("Error adding job:" + err.Error())
	}
	// 启动定时器
	crontab.Start()
	select {}
}

type myLogger struct {
}

func (ml myLogger) Info(msg string, keysAndValues ...interface{}) {
	zap.S().Info(fmt.Sprintf(msg, keysAndValues...))
}

func (ml myLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	zap.S().Error(err.Error() + fmt.Sprintf(msg, keysAndValues...))
}
