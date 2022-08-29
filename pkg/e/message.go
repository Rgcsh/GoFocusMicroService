package e

var Message = map[int]string {
	200: "Success",
	500: "Server Internal error",
	501: "DB error",
	502: "唯一性索引错误,数据插入重复",
	503: "数据库中找不到数据",
	504: "入参解析失败",
	505: "入参中必须只有一个值为列表",

}

func GetMessage(code int) string {
	message, ok := Message[code]
	if ok {
		return message
	}
	return Message[500]
}
