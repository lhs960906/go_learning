package function_define

import "fmt"

func Sum(a int, others ...int) int {
	sum := a
	for _, ele := range others {
		sum += ele
	}
	fmt.Printf("other: len %d, cap %d\n", len(others), cap(others))
	return sum
}
