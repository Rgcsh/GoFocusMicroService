package crontabs

import (
	"GoFocusMicroService/pkg/gcache"
	"GoFocusMicroService/pkg/utils"
	"fmt"
	"time"
)

func Job1() {
	fmt.Println("job1")
	err := gcache.Rds.Set("job1", "testRedis", time.Second*1)
	utils.PanicOnError(err, "操作redis失败")
	val, err := gcache.Rds.Get("job1")
	utils.PanicOnError(err, "操作redis失败")
	fmt.Printf("获取到redis key=job1,值为:%v", val)
}
