package utils

import (
	"encoding/json"
	"strings"
)

// RepeatChar 生成重复的字符
func RepeatChar(char byte, cnt int) string {
	s := make([]byte, cnt)
	for idx := 0; idx < cnt; idx++ {
		s[idx] = char
	}
	return string(s)
}

// FmtMobileStar 将手机号中间四位替换为*号
func FmtMobileStar(mobile string) string {
	if len(mobile) < 3 {
		return mobile
	}

	if len(mobile) < 7 {
		return mobile[:3] + RepeatChar('*', len(mobile)-3)
	}

	return mobile[:3] + RepeatChar('*', len(mobile)-7) + mobile[len(mobile)-4:]
}

// FmtEmailStar 将邮箱地址转换为带星号的模式
func FmtEmailStar(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}

	res := strings.Split(email, "@")
	if len(res[0]) == 1 {
		return "*@" + res[1]
	}
	if len(res[0]) == 2 {
		return res[0][:1] + "*@" + res[1]
	}
	return res[0][:2] + RepeatChar('*', len(res[0])-2) + "@" + res[1]
}

//
// SliceIn
//  @Description: 切片的in查询方式
//  @param slice:
//  @param num:
//  @return bool:
//
func SliceIn(slice []int, num int) bool {
	for _, sliceNum := range slice {
		if sliceNum == num {
			return true
		}
	}
	return false
}

//
//  @Description: 将字符串转为 切片数据,失败就返回nil
//  @param s:输入的字符串内存地址
//  @return []string: 字符串类型的切片内存地址
//
func JsonUnmarshal2Split(s *string) *[]string {
	var l []string
	err := json.Unmarshal([]byte(*s), &l)
	if err != nil {
		return nil
	}
	return &l
}

//
//  @Description: 检查入参 是否为 nil或"" ,如果是 返回 true
//  @param v:
//  @return bool:
//
func CheckValEmpty(v *interface{}) bool {
	if v == nil || *v == "" {
		return true
	}
	return false
}

//
//  @Description: 蛇形转为小驼峰
//  @param s:
//  @return *string:
//  snake2Camel("a_b")=="aB"
//
func Snake2Camel(s string) string {
	//切分,按照 _
	split := strings.Split(s, "_")
	//没有_,直接返回
	if len(split) == 1 {
		return s
	}
	//对第一个之后的每个数据尝试首字母转为大写,数字除外
	finalS := ""
	for index, subS := range split {
		if index == 0 {
			finalS += subS
		} else {
			finalS += FirstString2Upper(subS)
		}
	}
	return finalS
}

//
//  @Description: 字符串首字母大写
//  @param s:
//  @return string:
//  firstString2Upper("abc")=="Abc"
//
func FirstString2Upper(s string) string {
	if s == "" {
		return s
	}
	finalS := strings.ToUpper(s[:1]) + s[1:]
	return finalS
}
