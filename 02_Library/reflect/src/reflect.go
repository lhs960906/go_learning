package src

import (
	"fmt"
	"reflect"
)

func ReflectKind() {
	type MyInt int
	var x MyInt = 7
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())                        // src.MyInt
	fmt.Println("kind is int: ", v.Kind() == reflect.Int) // true.
}
