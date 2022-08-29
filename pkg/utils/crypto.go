package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

func MD5(val string) []byte {
	w := md5.New()
	_, _ = io.WriteString(w, val)
	return w.Sum(nil)
}

func Md5(val string) string {
	return hex.EncodeToString(MD5(val))
}

// Seed 用于生成指定长度的随机字符串，由 a-z, A-Z, 0-9 组成
// Seed func is used to generate a random string which is combined with a-z, A-Z, 0-9.
func Seed(size int) string {
	// 初始参数
	result, kindLen, kindStart := make([]byte, size), [3]int{10, 26, 26}, [3]int{48, 97, 65}

	// 修改随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成代码
	for i := 0; i < size; i++ {
		kind := rand.Intn(3)
		kLen, kStart := kindLen[kind], kindStart[kind]
		result[i] = uint8(kStart + rand.Intn(kLen))
	}
	return string(result)
}
