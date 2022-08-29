package focus

import (
	"GoFocusMicroService/models/focus_go"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/glog"
	"GoFocusMicroService/pkg/gtime"
	"GoFocusMicroService/pkg/utils"
	"github.com/gin-gonic/gin"
)

type QueryFields struct {
	Uid         int    `json:"uid" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	FocusedType string `json:"focused_type"`
	FocusedData string `json:"focused_data"`
	PageIndex   int    `json:"page_index" binding:"required"`
	PageSize    int    `json:"page_size" binding:"required"`
}

//
//  @Description: 查询 关注数据 或检查是否添加关注
//  四个参数都传 可以检查 对应数据是否添加关注,返回的切片长度为0 表示未关注
//  此接口可以灵活的查询不同范围的数据
//  @param c:
//
func QueryFocusHandler(c *gin.Context) {
	// 用户入参
	var params QueryFields
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)

	// 结果数据结构
	result := make(map[string]interface{})

	glog.Log.Info("查询分页数据")
	queryFocusParams := focus_go.Focus{
		Uid:         params.Uid,
		ServiceName: params.ServiceName,
		FocusedType: params.FocusedType,
		FocusedData: params.FocusedData}
	pageDbResult, err := focus_go.PagingQueryFocus(&queryFocusParams, params.PageIndex, params.PageSize)
	utils.ErrorCheck(err)
	pageResult := QueryFocusFormatData(pageDbResult)
	result["pageResult"] = pageResult

	glog.Log.Info("查询数据总个数")
	count, err := focus_go.QueryCountFocus(&queryFocusParams)
	utils.ErrorCheck(err)
	result["count"] = count

	ctx.Success(result)
}

//
// formatData
//  @Description: 对查询的结果 进行二次处理
//  @param params:
//  @param resultList:
//  @return map[string]interface{}:
//
func QueryFocusFormatData(resultList *[]focus_go.Focus) (result []interface{}) {
	for _, item := range *resultList {
		focus := make(map[string]interface{})
		focus["id"] = item.Id
		focus["createTime"] = gtime.FormatTimeLocal(item.PublicFields.CreateTime)
		focus["uid"] = item.Uid
		focus["serviceName"] = item.ServiceName
		focus["focusedType"] = item.FocusedType
		focus["focusedData"] = item.FocusedData
		result = append(result, focus)
	}
	return result
}
