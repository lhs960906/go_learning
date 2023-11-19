package function_define

import (
	"fmt"
	"testing"
)

func TestSum(testing *testing.T) {
	fmt.Println(Sum(0, 1))
	fmt.Println(Sum(0, 1, 2))
	fmt.Println(Sum(0, 1, 2, 3))
}
