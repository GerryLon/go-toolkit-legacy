// 控制同时执行的goroutine个数
// ref: https://github.com/EDDYCJY/gsema/blob/master/sema.go
package gsema

import "sync"

type GSemaphore struct {
	ch chan bool
	wg *sync.WaitGroup
}

func NewSemaphore(maxChan int) *GSemaphore {
	return &GSemaphore{
		ch: make(chan bool, maxChan),
		wg: &sync.WaitGroup{},
	}
}

// wg.Add()
func (gs *GSemaphore) Add(n int) {
	gs.wg.Add(n)
	for i := 0; i < n; i++ {
		gs.ch <- true
	}
}

// wg.Done()
// 表示一个goroutine完成
func (gs *GSemaphore) Done() {
	<-gs.ch
	gs.wg.Done()
}

// wg.Wait()
func (gs *GSemaphore) Wait() {
	gs.wg.Wait()
}
