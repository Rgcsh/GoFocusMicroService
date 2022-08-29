package routers

import (
	"GoFocusMicroService/conf"
	"GoFocusMicroService/controllers/focus"
	"GoFocusMicroService/controllers/rabbitmq_product"
	"GoFocusMicroService/crontabs"
	"GoFocusMicroService/middleware"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

//
// InitRouter
//  @Description: 注册 路由/定时任务
//  @return *gin.Engine:
//
func InitRouter() *gin.Engine {
	r := gin.New()
	// 是否启动 性能分析功能 true:启动
	if conf.Conf.App.EnablePProf {
		pprof.Register(r)
	}

	// 路由注册
	// 关注功能
	apiFocus := r.Group("/focus", middleware.CostTime(), gin.RecoveryWithWriter(nil, middleware.PanicRecovery))
	apiFocus.POST("/add", middleware.GetUid(), middleware.SynchronizedApi("GoFocusMicroService/focus/add"), focus.AddFocusHandler)
	apiFocus.POST("/del", focus.DelFocusHandler)
	apiFocus.POST("/query", focus.QueryFocusHandler)
	apiFocus.POST("/query/batch", focus.QueryBatchFocusHandler)

	// rabbitmq发送消息功能
	apiRabbitmq := r.Group("/rabbitmq", middleware.CostTime(), gin.RecoveryWithWriter(nil, middleware.PanicRecovery))
	apiRabbitmq.POST("/product", rabbitmq_product.RabbitmqProductHandler)

	//定时任务注册
	//每分钟的1s执行任务
	_, err := crontabs.Crontab.AddFunc("1 * * * * *", crontabs.Job1)
	if err != nil {
		panic(fmt.Sprintf("定时任务注册失败:%v", err))
	}
	return r
}
