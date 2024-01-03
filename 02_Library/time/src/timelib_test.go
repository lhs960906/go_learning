package time_lib

import (
	"testing"
)

// 获取1970到现在的 "秒/毫秒/微妙/纳秒" 数
func TestTimeStamp(t *testing.T) {
	TimeStamp()
}

// 获取时间的时间差
func TestTimeDelta(t *testing.T) {
	TimeDelta()
}

// 获取时间的年月日时分秒
func TestTimeEle(t *testing.T) {
	TimeEle()
}

// 时间格式化
func TestFormat(t *testing.T) {
	TimeFormat()
}

// 时间转 Time
func TestString2Time(t *testing.T) {
	String2Time()
}
