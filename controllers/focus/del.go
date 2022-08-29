package focus

import (
	"GoFocusMicroService/models/focus_go"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DelFields struct {
	Uid         int    `json:"uid" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	FocusedType string `json:"focused_type" binding:"required"`
	FocusedData string `json:"focused_data" binding:"required"`
}

//
// DelFocusHandler
//  @Description: 取消关注 物理删除关注数据
//  @param c:
//
func DelFocusHandler(c *gin.Context) {
	// 用户入参
	var params DelFields
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)

	err := focus_go.DeleteFocus(&params.Uid, &params.ServiceName, &params.FocusedType, &params.FocusedData)
	//  对异常进行 详细分类 处理;类似python的try catch
	utils.HandlerErrorAndResponse(ctx, err)
}
