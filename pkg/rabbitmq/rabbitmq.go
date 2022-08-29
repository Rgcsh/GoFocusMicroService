// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/7/27 15:57'
//
// Usage:
//
package rabbitmq

import (
	"GoFocusMicroService/conf"
	"GoFocusMicroService/pkg/utils"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

//rabbitmq的channel
type RabbitmqChannel struct {
	ch *amqp.Channel
}

var RabChan *RabbitmqChannel

//
//  RabbitmqQueue
//  @Description: rabbitmq消息队列,每次需要新建队列时,都要在下面声明
//
type RabbitmqQueue struct {
	Name  string
	ch    *amqp.Channel
	Queue *amqp.Queue
}

//发送短信的队列
var SmsRabQu *RabbitmqQueue

//发送邮件的队列
var EmailRabQu *RabbitmqQueue

//
//  @Description: rabbitmq初始化,只在项目启动时加载一次
//
func SetUp() {
	RabChan = NewRabbitmqClient(conf.Conf.RabbitmqConf)
	SmsRabQu = &RabbitmqQueue{Name: "GoFocusMicroService_sms", ch: RabChan.ch}
	EmailRabQu = &RabbitmqQueue{Name: "GoFocusMicroService_email", ch: RabChan.ch}
}

//
//  @Description: 新建 rabbitmq客户端,及channel
//  @param conf:
//  @return *RabbitmqChannel:
//
func NewRabbitmqClient(conf conf.RabbitmqConf) *RabbitmqChannel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", conf.UserName, conf.Password, conf.Host, conf.Port))
	if err != nil {
		utils.PanicOnError(err, "连接rabbitmq失败")
	}
	ch, err := conn.Channel()
	if err != nil {
		utils.PanicOnError(err, "新建rabbitmq channel失败")
	}
	return &RabbitmqChannel{ch: ch}
}

//
//  @Description: 声明一个队列(必须在发布和消费时都调用,因为rabbitmq删除队列后,重新尝试会报错 找不到 队列)
//  @receiver r
//  @param queueName:
//  @return *amqp.Queue:
//
func (r *RabbitmqQueue) QueueDeclare(queueName string) {
	//队列持久化
	queue, err := r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	utils.ErrorCheck(err)
	//在消费者处理完正在处理的任务之前,消费者不会预取新的任务;
	err = r.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.ErrorCheck(err)
	r.Queue = &queue
}

//
//  @Description: 发布消息到队列中
//  @receiver r
//  @param queue:
//  @param message:
//
func (r *RabbitmqQueue) Publish(message string) {
	r.QueueDeclare(r.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := r.ch.PublishWithContext(ctx, "", r.Queue.Name, false, false,
		// 消息持久化
		amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "application/json", Body: []byte(message)})
	utils.ErrorCheck(err)
}

//
//  @Description: 消费消息
//  @receiver r
//
func (r *RabbitmqQueue) Consume() (msgs <-chan amqp.Delivery) {
	r.QueueDeclare(r.Name)
	//关闭消息自动确认机制
	msgs, err := r.ch.Consume(r.Queue.Name, "", false, false, false, false, nil)
	utils.ErrorCheck(err)
	return msgs
}
