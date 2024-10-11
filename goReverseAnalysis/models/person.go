package models

import "fmt"

// Person 结构体定义
type Person struct {
	Name string
	Age  int
}

// NewPerson 构造函数
func NewPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

// Greet 方法
func (p *Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}
