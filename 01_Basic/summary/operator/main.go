package main

import "fmt"

func main() {
	fmt.Println(uint8(2) << 2) // 00000010 << 2 = 00001000；打印8
	fmt.Println(uint8(2) >> 2) // 00000010 << 2 = 00000000；打印0
}
