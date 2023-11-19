package files

import (
	"fmt"
	"io"
	"os"
)

func FileWrite(filepath string) error {
	fin, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file error! %s", err)
	}
	defer fin.Close()

	var n int
	if n, err = fin.Write([]byte("12345asdasfasgfasgagasg")); err != nil {
		return fmt.Errorf("write file error! %s", err)
	}
	fmt.Printf("write %d bytes to %s\n", n, filepath)

	return nil
}

func FileReadAll(filepath string) error {
	fin, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file error! %s", err)
	}
	defer fin.Close()

	b := make([]byte, 5)
	for {
		n, err := fin.Read(b)
		if err != nil {
			if err == io.EOF {
				// fmt.Println(n, string(b[:n]))
				break
			}
			return fmt.Errorf("read file error! %s", err)
		}
		fmt.Println(n, string(b[:n]))
	}

	return nil
}

func FileCopy(src string, dst string) error {
	// 只读打开源文件
	fsrc, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file %s error! err info: %s", src, err)
	}
	defer fsrc.Close()

	// 只写打开目标文件
	fdst, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("open file %s error! err info: %s", dst, err)
	}
	defer fdst.Close()

	// 将源文件内容拷贝到目标文件
	b := make([]byte, 10)
	for {
		n, err := fsrc.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read file err in copy! err info: %s", err)
		}
		fdst.Write(b[:n])
	}

	return nil
}
