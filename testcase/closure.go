/**
 * Introduction 匿名函数
 * Level 2 a + b 负样本
 */
// evaluation information start
// real case = false
// evaluation item = 完整度->链路完整度->函数调用->匿名函数和闭包->返回值传递
// bind_url = /completeness/chain_tracing/function_call/anonymous_function_closure/anonymous_function_001_F/anonymous_function_001_F.go
// evaluation information end

package main

import (
	"fmt"
	"time"
)

func anonymous_function_001_F(__taint_src string) {
	process := func(a string, b string) string {
		if a == "unsafe" {
			return a
		}
		return a + b
	}
	go process("a","b")
	__taint_sink(process(__taint_src, "foo"))
}

func __taint_sink(o interface{}) {
	fmt.Println(o)
}

func main() {
	go anonymous_function_001_F("unsafe")
	time.Sleep(1 * time.Second)
}
