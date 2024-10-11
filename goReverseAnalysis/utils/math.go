package utils

// Calculator 接口声明
type Calculator interface {
	Add(a, b int) int
	Subtract(a, b int) int
}

// BasicCalculator 实现 Calculator 接口
type BasicCalculator struct{}

func (c BasicCalculator) Add(a, b int) int {
	return a + b
}

func (c BasicCalculator) Subtract(a, b int) int {
	return a - b
}
