package rabbitmq_product

import (
	"GoFocusMicroService/controllers/public_method"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/rabbitmq"
	"GoFocusMicroService/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RabbitmqProductHandler(c *gin.Context) {
	// 用户入参
	var params interface{}
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)

	//异步发送短信
	smsTask := public_method.SmsTask{
		Account:       "17788888888",
		SignName:      "短信发送签名",
		TemplateCode:  "test_模板好",
		TemplateParam: fmt.Sprintf("{\"code\":\"%s\"}", "1234")}
	rabbitmq.SmsRabQu.Publish(utils.JsonMarshal(smsTask))

	//异步发送邮件
	emailTask := public_method.EmailTask{
		Account:   "213@qq.com",
		Subject:   "登录验证码",
		FromAlias: "xxx公司",
		HtmlBody:  fmt.Sprintf("<p>验证码为:%v</p>", "1234")}
	rabbitmq.EmailRabQu.Publish(utils.JsonMarshal(emailTask))

	ctx.SuccessMap()
}
