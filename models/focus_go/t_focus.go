package focus_go

import (
	"GoFocusMicroService/pkg/api_error"
	"fmt"
	"reflect"
	"time"
)

//
//  PublicQueryFields
//  @Description: 公共字段
//
type PublicFields struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement;not null;column:id;comment:主键ID;"`
	State      int       `json:"state" gorm:"default:1;column:state;comment:状态 1:正常 0:删除;"`
	CreateTime time.Time `json:"createTime" gorm:"default:CURRENT_TIMESTAMP;column:create_time;comment:创建时间;"`
	UpdateTime time.Time `json:"updateTime" gorm:"default:CURRENT_TIMESTAMP;column:update_time;comment:更新时间;"`
}

//
//  Focus
//  @Description: 关注model结构体
//
type Focus struct {
	PublicFields
	Uid         int    `json:"uid" gorm:"not null;column:uid;comment:用户ID;"`
	ServiceName string `json:"serviceName" gorm:"not null;column:service_name;comment:微服务项目名称(最好写对应项目名称);"`
	FocusedType string `json:"focusedType" gorm:"not null;column:focused_type;comment:被关注数据类型(最好写英文含义);"`
	FocusedData string `json:"focusedData" gorm:"not null;column:focused_data;comment:被关注的数据;"`
}

//
var PublicQueryFields = []string{"id", "state", "create_time", "update_time"}

// 查询的字段
var FocusQueryFields = []string{"id", "uid", "service_name", "focused_type", "focused_data", "create_time", "update_time"}

//表名称
//var tableName = "t_focus"

//
// InsertFocus
//  @Description: 新增一条关注数据
//  @param info:
//  @return err:
//
func InsertFocus(info *Focus) (err error) {
	return db.Omit(PublicQueryFields...).Create(&info).Error
}

//
// DeleteFocus
//  @Description: 删除关注数据
//  @param Id:主键ID的指针
//  @return err:
//
func DeleteFocus(Uid *int, ServiceName *string, FocusedType *string, FocusedData *string) (err error) {
	return db.Model(&Focus{}).Where(&Focus{Uid: *Uid, ServiceName: *ServiceName, FocusedType: *FocusedType, FocusedData: *FocusedData}).Delete(Focus{}).Error
}

//
// CheckExistsFocus
//  @Description: 检查关注是否存在
//  @param info:
//  @return f:
//  @return err:
//
func CheckExistsFocus(info *map[string]interface{}) (f *Focus, err error) {
	err = db.Model(&Focus{}).Where(&info).First(f).Error
	return
}

//
// PagingQueryFocus
//  @Description: 分页查询关注的数据
//  @param info:被查询的条件
//  @param pageIndex:
//  @param pageSize:
//  @return *[]Focus:
//  @return error:
//
func PagingQueryFocus(info *Focus, pageIndex int, pageSize int) (f *[]Focus, err error) {
	err = db.Model(&Focus{}).Select(FocusQueryFields).Where(&info).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&f).Error
	return
}

//
// QueryCountFocus
//  @Description: 查询符合条件的数据个数
//  @param info:
//  @return count:
//  @return err:
//
func QueryCountFocus(info *Focus) (count int64, err error) {
	err = db.Model(&Focus{}).Where(&info).Count(&count).Error
	return
}

//
//  @Description: 批量(In)查询数据
//  @param Uid:类型为 int或[]float64
//  @param ServiceName:类型为 string或[]string
//  @param FocusedType:类型为 string或[]string
//  @param FocusedData:类型为 string或[]string
//  @return f:
//  @return sliceParams:入参为切片类型的字段
//  @return err:
//
func QueryBatchFocus(Uid, ServiceName, FocusedType, FocusedData interface{}) (f *[]Focus, sliceParams []string, err error) {
	// in查询的参数,待会需要
	sql := db.Model(&Focus{}).Select(FocusQueryFields).Where("state=1")
	for _, item := range [4][2]interface{}{
		{Uid, "uid"},
		{ServiceName, "service_name"},
		{FocusedType, "focused_type"},
		{FocusedData, "focused_data"},
	} {
		obj := item[0]
		field := item[1]
		//入参尝试转换类型,在生成对应的sql
		switch reflect.ValueOf(obj).Kind() {
		case reflect.Float64:
			//	转为int类型,再 sql = 操作
			intObj := int(obj.(float64))
			sql.Where(fmt.Sprintf("%v=?", field), intObj)
		case reflect.String:
			//	判断是否为空字符串,不为空 就= 操作
			strObj := obj.(string)
			if strObj != "" {
				sql.Where(fmt.Sprintf("%v=?", field), strObj)
			}
		case reflect.Slice:
			//  直接sql in 操作
			sliceParams = append(sliceParams, field.(string))
			obj := obj.([]interface{})
			var split []interface{}
			// 将切片里的数据 转为int类型或string类型
			for _, sliceItem := range obj {
				if reflect.ValueOf(sliceItem).Kind() == reflect.Float64 {
					sliceItemObj := int(sliceItem.(float64))
					split = append(split, sliceItemObj)
				} else {
					sliceItemObj := sliceItem.(string)
					split = append(split, sliceItemObj)
				}
			}
			sql.Where(fmt.Sprintf("%v in ?", field), split)
		case reflect.Invalid:
			//  直接跳过
		}
	}

	if len(sliceParams) != 1 {
		return nil, nil, api_error.New(505)
	}
	err = sql.Debug().Find(&f).Error
	return
}
