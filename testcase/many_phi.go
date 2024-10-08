package main

import (
	"fmt"
	"math/rand"
)

func main() {
	a := 3
	b := 4
	c := 12
	x := 114
	y := 514
	if rand.Int()%2 == 0 {
		a = x + y
		b = a + x
	} else {
		a = b + 2
		c = y + 1
	}
	a = c + a
	fmt.Println(a)
}
