package main

import (
	"fmt"
	"log"

	"goReverseAnalysis/concurrency"
	"goReverseAnalysis/models"
	"goReverseAnalysis/utils"
)

func main() {
	// 使用结构体和方法
	person := models.NewPerson("Alice", 30)
	fmt.Println(person.Greet())

	// 使用接口
	var calculator utils.Calculator = utils.BasicCalculator{}
	result := calculator.Add(10, 5)
	fmt.Printf("10 + 5 = %d\n", result)

	// 使用goroutine和channel
	concurrency.RunWorkers()

	// 错误处理示例
	if err := doSomethingRisky(); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	fmt.Println("All tasks completed successfully.")
}

func doSomethingRisky() error {
	return fmt.Errorf("something went wrong")
}
