package example

import (
	"fmt"
	"math/rand"
	"time"
)

type Work struct {
	ID      int
	Content int64
}

var workChannel = make(chan *Work)

func doWork(id int) {
	fmt.Printf("Starting work: %d\n", id)

	random := rand.Int63n(10)
	time.Sleep(time.Duration(random) * time.Second)

	result := &Work{
		ID:      id,
		Content: random,
	}

	workChannel <- result
}

func main() {
	fmt.Println("ETL process started")

	rand.Seed(time.Now().UnixNano()) // Random seed 설정은 main에서 1회만

	start := time.Now()

	for i := 0; i < 10; i++ {
		go doWork(i)
	}

	for i := 0; i < 10; i++ {
		result := <-workChannel
		fmt.Printf("Work %d completed with result: %d\n", result.ID, result.Content)
	}

	duration := time.Since(start)
	fmt.Printf("ETL process completed in %.2f seconds\n", duration.Seconds())
}
