package main

import "fmt"

func main() {
	cond := true
	var o interface{}
	if cond {
		o = 1
	} else {
		o = "hello"
	}
	s := o.(string)
	fmt.Println(s)
}
