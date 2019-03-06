package main

import (
	"fmt"
	"github.com/GerryLon/go-toolkit/gsema"
	"time"
)

var sema = gsema.NewSemaphore(3)

// output
// seq: 9, time: 1550465827
// seq: 3, time: 1550465827
// seq: 2, time: 1550465827
// seq: 4, time: 1550465828
// seq: 5, time: 1550465828
// seq: 6, time: 1550465828
// seq: 8, time: 1550465829
// seq: 7, time: 1550465829
// seq: 0, time: 1550465829
// seq: 1, time: 1550465830
//
// Process finished with exit code 0
func main() {
	userCount := 10 // 需要执行的goroutine数

	for i := 0; i < userCount; i++ {
		go Read(i)
	}
	sema.Wait()
}

func Read(i int) {
	defer sema.Done()
	sema.Add(1)
	fmt.Printf("seq: %d, time: %d\n", i, time.Now().Unix())
	time.Sleep(time.Second)
}
