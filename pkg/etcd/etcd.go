package etcd

import (
	"GoFocusMicroService/conf"
	"GoFocusMicroService/pkg/api_error"
	"GoFocusMicroService/pkg/glog"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

var EtcdCli *clientv3.Client

// SetUp 实例化etcd
func SetUp() {
	etcdCli, err := clientv3.New(clientv3.Config{Endpoints: []string{conf.Conf.Etcd.GrpcProxy}, DialTimeout: 5 * time.Second})
	if err != nil {
		panic(fmt.Sprintf("连接etcd失败:%v", err.Error()))
	}
	EtcdCli = etcdCli
}

func SynchronizedApiByEtcd(name string, expire int) {
	//新建一个lease(租约)
	//注意 此租约会一直刷新超时时间,保证未手动关闭租约时此租约一直有效,此方式能保证 在后面的获取锁,及释放锁 时间段内 锁能一直有效,最大化防止并发错误;
	//租约退出条件为:直到 手动关闭租约或程序异常退出,租约里的锁会在声明的expire时间后被删除,保证不会死锁;
	session, err := concurrency.NewSession(EtcdCli, concurrency.WithTTL(expire))
	if err != nil {
		glog.Log.Error("获取lease(租约)失败")
		panic(err)
	}
	//最后关闭lease
	defer func() { _ = session.Close() }()

	//新建一个锁对象
	l := concurrency.NewMutex(session, name)

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

	//todo:进行业务层操作
	time.Sleep(8 * time.Second)

	//开始释放锁
	err = l.Unlock(context.TODO())
	if err != nil {
		glog.Log.Info("删除锁对应的key失败,后续etcd应会自动ttl删除")
		panic(err)
	}
	glog.Log.Info("释放锁成功")
}
