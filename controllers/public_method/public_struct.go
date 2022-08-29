package public_method

//
//  @Description: 短信发送任务
//
type SmsTask struct {
	Account       string //发送目标手机号
	SignName      string //短信签名
	TemplateCode  string //短信模板代码
	TemplateParam string //模板参数
}

//
//  @Description: 邮件发送任务
//
type EmailTask struct {
	Account   string //发送目标手机号
	Subject   string //邮件主题
	FromAlias string //邮件...
	HtmlBody  string //邮件内容
}
