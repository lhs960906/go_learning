package files_test

import (
	"testing"

	"github.com/lhs960906/Go-Learning/01_Basic/io/files"
)

// 测试文件写入
func TestFileWrite(t *testing.T) {
	files.FileWrite("test.txt")
}

// 测试读取文件的全部内容
func TestReadFileAll(t *testing.T) {
	files.FileReadAll("test.txt")
}

// 测试文件复制
func TestFileCopy(t *testing.T) {
	files.FileCopy("林乐清.png", "林乐清_c.png")
}
