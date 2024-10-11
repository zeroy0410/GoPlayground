package concurrency

import (
	"fmt"
	"time"
)

// RunWorkers 展示了并发和通道的使用
func RunWorkers() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for r := 1; r <= 5; r++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}
