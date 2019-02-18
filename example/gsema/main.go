package main

import (
	"fmt"
	"github.com/GerryLon/go-toolkit/gsema"
	"time"
)

var sema = gsema.NewSemaphore(3)

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
