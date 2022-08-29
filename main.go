package main

import (
	"GoFocusMicroService/conf"
	"GoFocusMicroService/controllers/consumers"
	"GoFocusMicroService/crontabs"
	"GoFocusMicroService/models"
	"GoFocusMicroService/pkg/etcd"
	"GoFocusMicroService/pkg/gcache"
	"GoFocusMicroService/pkg/glog"
	"GoFocusMicroService/pkg/rabbitmq"
	"GoFocusMicroService/pkg/validator_rewrite"
	"GoFocusMicroService/routers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	conf.SetUp()
	crontabs.SetUp()
	glog.SetUp()
	gcache.SetUp()
	rabbitmq.SetUp()
	models.SetUp()
	consumers.SetUp()
	validator_rewrite.SetUp()
	etcd.SetUp()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	// 设置未声明的参数无法传参
	gin.EnableJsonDecoderDisallowUnknownFields()
	r := routers.InitRouter()

	glog.Log.Info("http server will start at", zap.Int("port", 7066))
	if err := r.Run(":7066"); err != nil {
		glog.Log.Error("http server start failed", zap.Error(err))
		return
	}
}
