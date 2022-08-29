// Package middleware
// @Description: 中间件 可以注册到 全局路由,或者某个路由组 或者单个路由中
//方法如下链接的 注册中间件 部分:https://www.liwenzhou.com/posts/Go/Gin_framework/
package middleware

import (
	"GoFocusMicroService/pkg/api_error"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/e"
	"GoFocusMicroService/pkg/etcd"
	"GoFocusMicroService/pkg/glog"
	"GoFocusMicroService/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

//
// CostTime
//  @Description: 定义一个 统计接口耗时的中间件
//  @return gin.HandlerFunc:
//
func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		//继续执行请求的后续代码
		c.Next()
		// 计算距离当前时间的耗时
		costTime := time.Since(startTime)
		glog.Log.Info(fmt.Sprintf("接口:%v,耗时:%v", c.Request.URL, costTime))
	}
}

//
//  @Description: 接口级别的 分布式互斥锁;实现同一时间只能有一个服务提供访问,防并发操作
//  @param name: 锁的名字
//  @return gin.HandlerFunc:
//
func SynchronizedApi(name string) gin.HandlerFunc {

	var l *concurrency.Mutex

	synchronizedApiByEtcdLock := func(name string) {
		//新建一个lease(租约)
		//注意 此租约会一直刷新超时时间(自动续约),保证未手动关闭租约时此租约一直有效,此方式能保证 在后面的获取锁,及释放锁 时间段内 锁能一直有效,最大化防止并发错误;
		//所以 不需要用户设置过期时间,只需要写死1s保证最小超时时间即可;这样即使 程序异常退出,1s后可以获取锁,最小化减少锁的影响时间;
		//租约退出条件为:直到 手动关闭租约或程序异常退出,租约里的锁会在声明的expire时间后被删除,保证不会死锁;
		session, err := concurrency.NewSession(etcd.EtcdCli, concurrency.WithTTL(1))
		if err != nil {
			glog.Log.Error("获取lease(租约)失败")
			panic(err)
		}
		//最后关闭lease
		defer func() { _ = session.Close() }()

		//新建一个锁对象
		l = concurrency.NewMutex(session, name)

		//设置100毫秒后不再尝试获取锁
		c1 := context.Background()
		c2, cancel := context.WithTimeout(c1, 100*time.Millisecond)
		defer cancel()

		//尝试上锁
		err = l.Lock(c2)
		if err != nil {
			glog.Log.Info("上锁失败")
			panic(api_error.New(504, string("接口不能并发请求")))
		}
		glog.Log.Info("上锁成功")
	}

	synchronizedApiByEtcdUnLock := func() {
		//开始释放锁
		err := l.Unlock(context.TODO())
		if err != nil {
			glog.Log.Info("删除锁对应的key失败,后续etcd应会自动ttl删除")
			panic(err)
		}
		glog.Log.Info("释放锁成功")
	}

	return func(c *gin.Context) {
		// 获取uid
		uid, exists := c.Get("uid")
		if !exists {
			// 如果获取不到用户ID(Access-User),则此中间件会让所有用户的只要调用对应的接口就会防并发,不符合需求;
			glog.Log.Warn("必须在请求头中设置Access-User的值,且路由注册时先注册middleware.GetUid(),才能正常使用此中间件,此次调用 防并发中间件不生效!")
			c.Next()
		} else {
			//生成锁的key
			key := fmt.Sprintf("%v:%v", name, uid)
			glog.Log.Info(fmt.Sprintf("上锁key:%v", key))

			//上锁
			synchronizedApiByEtcdLock(key)
			//继续执行请求的后续代码
			c.Next()
			//释放锁
			synchronizedApiByEtcdUnLock()
		}
	}
}

//
// PanicRecovery
//  @Description: 此处捕捉到错误后,判断错误是否为 可以处理的错误,如果是就 返回相应的错误码给前端
//  此方法给 CustomRecovery 中间件使用
//  @param c:
//  @param err:
//
func PanicRecovery(c *gin.Context, err interface{}) {
	code, message := utils.ModelErrorHandler(err)
	if message == "" {
		c.AbortWithStatusJSON(200, app.NewResponse(code, e.GetMessage(code), make(map[string]interface{})))
	} else {
		c.AbortWithStatusJSON(200, app.NewResponse(code, message, make(map[string]interface{})))
	}
}

//
//  @Description: 从请求头中获取uid(用户ID)
//  @return gin.HandlerFunc:
//
func GetUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetHeader("Access-User")
		if uid != "" {
			c.Set("uid", uid)
		}
		c.Next()
	}
}
