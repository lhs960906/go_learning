package time_lib

import (
	"fmt"
	"time"
)

var location *time.Location

func init() {
	// 获取指定时区
	location, _ = time.LoadLocation("Asia/Shanghai")
}

// 获取 1970-01-01 00:00:00 到现在的 "秒/毫秒/微妙/纳秒" 数
func TimeStamp() {
	// 返回当前的本地时间
	t := time.Now()
	// 返回不同单位的时间戳
	fmt.Printf("TimeStamp1: %v\n", t.Unix())      // 秒 10^0(E0)
	fmt.Printf("TimeStamp2: %v\n", t.UnixMilli()) // 毫秒 10^-3(E-3)
	fmt.Printf("TimeStamp3: %v\n", t.UnixMicro()) // 微秒 10^-6(E-6)
	fmt.Printf("TimeStamp4: %v\n", t.UnixNano())  // 纳秒 10^-9(E-9)
}

// 时间加减操作
func TimeDelta() {
	// 返回当前的本地时间
	begin := time.Now()
	// 经过一段时间
	time.Sleep(1 * time.Second)

	// 获取当前时间到begin时间的差值
	timeDelta := time.Since(begin)
	fmt.Printf("TimeDelta1: %v\n", timeDelta)

	// 获取end时间到begin时间段的差值
	end := time.Now()
	timeDelta = end.Sub(begin)
	fmt.Printf("TimeDelta2: %v\n", timeDelta)

	// 获取begin时间8个小时后的时间
	eightHour := time.Duration(8 * time.Hour)
	end = begin.Add(eightHour)
	fmt.Printf("TimeDelta3: %v\n", end)
}

// 获取对当前的年月日时分秒
func TimeEle() {
	now := time.Now()
	// 年
	fmt.Printf("TimeEle1: %v\n", now.Year())
	// 月
	fmt.Printf("TimeEle2: %v\n", now.Month())
	// 日
	fmt.Printf("TimeEle3: %v\n", now.Day())
	// 时
	fmt.Printf("TimeEle4: %v\n", now.Hour())
	// 分
	fmt.Printf("TimeEle5: %v\n", now.Minute())
	// 秒
	fmt.Printf("TimeEle6: %v\n", now.Second())
}

// 时间格式化
func TimeFormat() {
	now := time.Now()
	// go 语言中格式化时间为 "yyyy-MM-dd hh:mm:ss" 形式, 必须使用 "2006-01-02 15:04:05"
	layout := "2006-01-02 15:04:05"
	fmt.Printf("TimeFormat: %v\n", now.Format(layout))
}

// 字符串转time.Time
func String2Time() {
	const (
		TIME_FMT1 = "2006-01-02 15:04:05"
		TIME_FMT2 = "2006-01-02"
		TIME_FMT3 = "20060102"
	)

	// 将字符串转换为Time(不推荐使用)
	if t, err := time.Parse(TIME_FMT2, "1992-02-18"); err == nil {
		fmt.Printf("String2Time: %v\n", t)
	}

	// 将字符串转换为Time(推荐使用)
	if t, err := time.ParseInLocation(TIME_FMT2, "1992-02-18", location); err == nil {
		fmt.Printf("String2Time: %v\n", t)
	}
}
