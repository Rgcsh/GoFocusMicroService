# GoFocusMicroService

Go语言编写的关注功能的微服务

跨系统的关注功能;所有系统的关注都可以在此系统完成

借此功能搭建Go语言gin的项目框架,及相关实用功能,后续其他项目可以按照此项目框架进行搭建;

# 简要系统设计文档

[关注微服务开发设计文档.md](./docs/关注微服务开发设计文档.md)

# 安装依赖包

* 在项目根目录执行 如下命令设置代理

```shell
export GOPROXY=https://goproxy.cn/,https://mirrors.aliyun.com/goproxy/,direct
```

* 执行 ```go mod tidy``` 命令(会自动分析你当前项目所需要的依赖，并且将他们下载下来)
* 执行 ```go mod vendor``` 命令 将公共依赖包复制到此项目中(注意只赋值使用到的模块,使用新模块时需要从新赋值,不建议使用此功能),这时 只供自己项目使用

# 项目运行方法
* 修改 conf/test.yml中的配置;
* export LOC_CFG=/Users/xxx/GoFocusMicroService/conf/test.yml
* go run main.go

# 项目功能用法

### 实现缓存功能

* 使用redis进行缓存

### 实现数据持久化存储功能

* 使用mysql存储数据

### 中断请求功能

* 调用方法:在 API处理函数的任何地方,通过调用 ```panic(&api_err.ApiError{500})``` 实现 终止此API的后续处理,立刻返回相应错误码的响应;
* 实现原理:自定义error错误结构体,panic此错误结构体指针 后, recovery中间件尝试捕捉此错误,并返回给前端相应结果;

### 性能分析启动及操作方式

* 启动方式:配置中 App.EnablePProf 设置为 true
* 分析操作方式: [性能分析方法](http://liumurong.org/2019/12/gin_pprof/)

### 添加定时任务方法

* 在crontabs包中添加定时任务,然后在 routers包中注册定时任务

### 每个用户接口层防并发(etcd分布式锁实现)功能

* 利用etcd自带的分布式锁功能实现接口及防并发功能; 路由注册时调用 middleware.SynchronizedApi() 中间件方法即可;

### 异步任务功能(rabbitmq生产/消费)

* 利用rabbitmq实现异步任务的生产和消费,如发送短信/邮件等功能

### 对入参校验的返回值改为中文

* gin默认的参数校验功能使用的是 github.com/go-playground/validator/v10 包,因此根据此包方法修改返回值为中文格式即可;

### 中间件(middleware)相关功能

* 请求耗时统计功能
* 接口层防并发功能
* 获取用户ID功能,且可以供后续handler使用
