package main

import (
	"bufio"
	"fmt"
	"github.com/GerryLon/go-toolkit/external_sort"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	// testRandomSource()
	testReaderSource()
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	max := 10 * count
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Intn(max)
		}
		close(out)
	}()
	return out
}

type ReaderSourceImpl struct {
}

func (r ReaderSourceImpl) Read() (<-chan int, error) {
	out := make(chan int)

	const testFile = "./numbers.txt"
	var err error
	var file *os.File
	var line []byte
	var number int
	var reader *bufio.Reader

	file, err = os.Open(testFile)
	if err != nil {
		close(out)
		return out, err
	}
	// 这里不能关闭, 不然goroutine中提示错误(文件已经关闭)
	// 这里暂时不做处理
	// defer file.Close()

	reader = bufio.NewReader(file)
	go func() {
		for {
			line, _, err = reader.ReadLine()
			if len(line) > 0 {
				number, err = strconv.Atoi(string(line))
				out <- number
			}

			if err != nil {
				break
			}
		}
		close(out)
	}()
	return out, err
}

func testReaderSource() {
	const testFile = "./numbers.txt"
	file, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入20个随机数
	for i := 0; i < 20; i++ {
		file.WriteString(strconv.Itoa(rand.Intn(100)) + "\n")
	}

	r := ReaderSourceImpl{}
	ch, err := external_sort.ReaderSource(r)
	if err != nil {
		panic(err)
	}
	ch = external_sort.MemorySort(ch)
	printChannel(ch)
}

func testRandomSource() {
	const channelCount = 10
	chs := make([]<-chan int, channelCount)
	for i := 0; i < channelCount; i++ {
		chs[i] = external_sort.MemorySort(
			RandomSource(10))
	}
	ch := external_sort.MergeN(chs...)
	printChannel(ch)
}

func printChannel(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
