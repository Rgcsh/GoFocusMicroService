package focus

import (
	"GoFocusMicroService/models/focus_go"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AddFields struct {
	Uid         int    `json:"uid" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	FocusedType string `json:"focused_type" binding:"required"`
	FocusedData string `json:"focused_data" binding:"required"`
}

func AddFocusHandler(c *gin.Context) {
	// 用户入参
	var params AddFields
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)

	err := focus_go.InsertFocus(&focus_go.Focus{Uid: params.Uid,
		ServiceName: params.ServiceName, FocusedType: params.FocusedType, FocusedData: params.FocusedData})
	//  对异常进行 详细分类 处理;类似python的try catch
	utils.HandlerErrorAndResponse(ctx, err)
}
