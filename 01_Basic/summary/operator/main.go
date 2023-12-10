package main

import "fmt"

func main() {
	fmt.Println(uint8(1) & uint8(1)) // 00000001 & 00000001 = 00000001；打印1
	fmt.Println(uint8(2) & uint8(1)) // 00000010 & 00000001 = 00000000；打印0
}
