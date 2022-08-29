// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/8/29 09:55'
//
// Usage:
//

package utils

import "encoding/json"

//
//  @Description: 将结构体或map转为json字符串
//  @param obj:
//  @return string:
//
func JsonMarshal(obj interface{}) string {
	sub, err := json.Marshal(obj)
	ErrorCheck(err)
	return string(sub)
}

//
//  @Description: 将json字符串改为结构体
//  @param str:
//  @param obj:
//
func JsonUnmarshal(str []byte, obj interface{}) {
	err := json.Unmarshal(str, obj)
	ErrorCheck(err)
	return
}
