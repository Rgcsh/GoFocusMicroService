package consumers

import (
	"GoFocusMicroService/controllers/public_method"
	"GoFocusMicroService/pkg/glog"
	"GoFocusMicroService/pkg/rabbitmq"
	"GoFocusMicroService/pkg/utils"
	"fmt"
)

//
//
//  @Description:邮件发送消费者
//  @param c:
//
// queue: GoFocusMicroService_email
func QueueConsumerEmail() {
	for {
		msgs := rabbitmq.EmailRabQu.Consume()
		for msg := range msgs {
			body := string(msg.Body)
			glog.Log.Info(fmt.Sprintf("邮件消费者接收到消息:%v", body))
			var emailTask public_method.EmailTask
			utils.JsonUnmarshal(msg.Body, &emailTask)
			glog.Log.Info(fmt.Sprintf("邮件发送任务结构体:%v", body))
			glog.Log.Info(fmt.Sprintf("邮件消费者处理完成:%v", body))

			//消息确认
			err := msg.Ack(false)
			if err != nil {
				glog.Log.Info(fmt.Sprintf("消息确认失败:%v", err))
			}
		}
		glog.Log.Info("队列可能出现被删除等情况,尝试重新连接...")
	}
}
