package gtime

import (
	"database/sql/driver"
	"fmt"
	"log"
	"time"
)

type Time time.Time

const (
	TimeFormat = "2006-01-02 15:04:05"
	zone       = "Asia/Shanghai"
)

// UnmarshalJson implements json unmarshal interface.
func (t *Time) UnmarshalJson(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

// MarshalJSON implements json marshal interface.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}

func (t Time) Local() time.Time {
	loc, _ := time.LoadLocation(zone)
	return time.Time(t).In(loc)
}

// Value ...
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// Scan valueof time.Time 注意是指针类型 method
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

//
// ParseLocal
//  @Description: 修改时间的时区为上海时区,且返回值为字符串格式
//  @param inputTime:
//  @return string:
//
func ParseLocal(inputTime time.Time) string {
	loc, _ := time.LoadLocation(zone)
	return inputTime.In(loc).Format(TimeFormat)
}

//
// FormatTimeLocal
//  @Description: 将时间修改为上海时区,并 返回 TimeFormat 字符串格式
//  @param inputTime:
//  @return string:
//
func FormatTimeLocal(inputTime time.Time) string {
	return inputTime.Format(TimeFormat)
}

//
// StringToTime
//  @Description: 字符串格式的时间 转为 Time对象的 时间
//  @param s:
//  @return time.Time:
//  @return error:
//
func StringToTime(s string, format string) (*time.Time, error) {
	log.Println()
	// 设置时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	timeObj, err := time.ParseInLocation(format, s, location)
	return &timeObj, err
}
