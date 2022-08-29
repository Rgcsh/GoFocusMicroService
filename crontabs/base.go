package crontabs

import (
	"github.com/robfig/cron/v3"
	"time"
)

var Crontab *cron.Cron

//
// SetUp
//  @Description: 定时任务初始化
//
func SetUp() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	Crontab = cron.New(cron.WithSeconds(), cron.WithLocation(location))
	//定时任务启动
	Crontab.Start()
}
