package main

import "fmt"

func complexFunction(x int, y int) int {
	// Variable declaration
	z := x * y

	// For loop with an if-else statement
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			z += i
		} else {
			z -= i
		}
	}

	// Switch statement
	switch z % 3 {
	case 0:
		fmt.Println("z is divisible by 3")
	case 1:
		fmt.Println("z has a remainder of 1 when divided by 3")
	default:
		fmt.Println("z has a remainder of 2 when divided by 3")
	}

	// A labeled loop with break and continue
OuterLoop:
	for j := 0; j < 3; j++ {
		for k := 0; k < 3; k++ {
			if j+k == 3 {
				break OuterLoop
			}
			if j == k {
				continue
			}
			z += j * k
		}
	}

	// Deferred function call
	defer fmt.Println("Function execution completed")

	// Anonymous function
	func() {
		fmt.Println("Anonymous function executed")
	}()

	// Nested if statement
	if z > 10 {
		if x > y {
			z += x
		} else {
			z -= y
		}
	} else {
		z = 10
	}

	return z
}

func main() {
	result := complexFunction(2, 3)
	fmt.Println("Result:", result)
}
