package main

import "fmt"

func main() {
	var i interface{} = "hello"

	// 错误的类型断言，i 实际上是 string 类型
	n := i.(int) // 这里会导致 panic
	fmt.Println(n)
}
