package focus

import (
	"GoFocusMicroService/models/focus_go"
	"GoFocusMicroService/pkg/app"
	"GoFocusMicroService/pkg/glog"
	"GoFocusMicroService/pkg/gtime"
	"GoFocusMicroService/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

//
//  QueryBatchFields
//  @Description: 接口入参限制,所有参数非必传字段,如果参数值为切片,就in条件查询sql, string/int类型就=条件查询sql
//
type QueryBatchFields struct {
	// 类型可以为 int,[]int
	Uid interface{} `json:"uid"`
	// 类型可以为 string,[]string
	ServiceName interface{} `json:"service_name"`
	// 类型可以为 string,[]string
	FocusedType interface{} `json:"focused_type"`
	// 类型可以为 string,[]string
	FocusedData interface{} `json:"focused_data"`
}

//
//  @Description: 批量查询 关注数据
//  此接口可以灵活的查询不同范围的数据
//  @param c:
//
func QueryBatchFocusHandler(c *gin.Context) {
	// 用户入参
	var params QueryBatchFields
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)
	fmt.Println(ctx)

	glog.Log.Info("查询分页数据")
	DbResult, sliceParams, err := focus_go.QueryBatchFocus(params.Uid, params.ServiceName, params.FocusedType, params.FocusedData)
	sliceParam := sliceParams[0]
	fmt.Println(sliceParams)
	utils.ErrorCheck(err)
	pageResult := GroupFocusData(DbResult, sliceParam)
	ctx.Success(pageResult)
}

//
//  @Description: 对查询结果根据 sliceParam进行分组
//  @param resultList:
//  @param sliceParam:
//  @return []interface{}:
//
func GroupFocusData(resultList *[]focus_go.Focus, sliceParam string) *[]map[string]interface{} {
	// {
	//  "uid1":[{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx}],
	//  "uid2":[{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx}]
	// }
	var mapObj map[interface{}][]interface{}
	mapObj = make(map[interface{}][]interface{})
	for _, item := range *resultList {
		focus := make(map[string]interface{})
		focus["id"] = item.Id
		focus["createTime"] = gtime.FormatTimeLocal(item.PublicFields.CreateTime)
		focus["uid"] = item.Uid
		focus["serviceName"] = item.ServiceName
		focus["focusedType"] = item.FocusedType
		focus["focusedData"] = item.FocusedData
		//根据分组key获取对应的值
		key := focus[utils.Snake2Camel(sliceParam)]
		//判断map中key是否存在
		if _, ok := mapObj[key]; ok {
			//	存在就append
			mapObj[key] = append(mapObj[key], focus)
		} else {
			//	不存在就新增
			mapObj[key] = []interface{}{focus}
		}
	}

	// [
	//	{"key":"uid1","focusList":[{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx}]},
	//	{"key":"uid2","focusList":[{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx},{xxx:xxx,xxx:xxx}]},
	//	]
	var SliceObj []map[string]interface{}
	//转换数据结构
	for key, Slice := range mapObj {
		item := map[string]interface{}{
			"key":       key,
			"focusList": Slice,
		}
		SliceObj = append(SliceObj, item)
	}
	return &SliceObj
}
