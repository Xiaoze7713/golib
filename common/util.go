/**
 * @Author: wenliangzhang
 * @Description:
 * @File: util.go
 * @Version: 1.0.0
 * @Date: 2022/5/24 2:58 PM
 */
package common

import (
	"time"
)

// UTimeMs 毫秒时间戳
func UTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetTimeDifferenceLastDay 获取当前时间到明天0点整的时间差
func GetTimeDifferenceLastDay() int64 {
	nowTime := time.Now()
	// 当天秒级时间戳
	nowTimeStamp := nowTime.Unix()

	nowTimeStr := nowTime.Format("2006-01-02")

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	// 第二天零点时间戳
	towTimeStamp := t2.AddDate(0, 0, 1).Unix()

	return towTimeStamp - nowTimeStamp
}
