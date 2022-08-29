package api_error

import "fmt"

//
//  ApiError
//  @Description: 自定义的错误结构体
//
type ApiError struct {
	ErrCode    int
	ErrMessage string
}

//
//  @Description:新建一个错误结构体对象
//  @param errs: 可变入参方式,第一个参数 为 int类型,表示错误码  第二个参数可选,为string类型,表示错误消息
//  @return error:
//
func New(errs ...interface{}) error {
	if len(errs) == 0 {
		panic("New参数至少存在一个")
	}
	errCode := errs[0].(int)
	var errMessage string
	if len(errs) == 2 {
		errMessage = errs[1].(string)
	} else {
		errMessage = ""
	}
	return &ApiError{errCode, errMessage}
}

//
// Error
//  @Description:
//  @receiver a
//  @return string:
//
func (a *ApiError) Error() string {
	return fmt.Sprintf("错误码:%v,错误消息:%v", a.ErrCode, a.ErrMessage)
}
