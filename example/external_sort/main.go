package main

import (
	"bufio"
	"fmt"
	"github.com/GerryLon/go-toolkit/external_sort"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
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

func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
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

	// 写入10个随机数
	for i := 0; i < 10; i++ {
		file.WriteString(strconv.Itoa(rand.Intn(100)) + "\n")
	}

	r := ReaderSourceImpl{}
	ch, err := external_sort.ReaderSource(r)
	if err != nil {
		panic(err)
	}

	startIntServer()
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	nig := NetIntGetter{
		conn: conn,
	}
	chx, err := nig.Read()
	if err != nil {
		panic(err)
	}

	// 将不同数据源的数据合并
	ch = external_sort.MergeN(
		external_sort.MemorySort(ch),
		external_sort.MemorySort(RandomSource(10)),
		external_sort.MemorySort(ArraySource(3, 5, 8, 10, 99, 6, 9999, 1023, 888, 250)),
		external_sort.MemorySort(chx),
	)
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
	i := 0
	for v := range ch {
		i++
		fmt.Println(i, v)
	}
}

type NetIntGetter struct {
	conn net.Conn
}

// 从网络上读取整形, 构建channel
func (n NetIntGetter) Read() (<-chan int, error) {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 16)
		strs := make([]string, 0)
		defer n.conn.Close()
		for {
			n, err := n.conn.Read(buffer)
			if n > 0 {
				strs = append(strs, string(buffer[:n]))
			} else {
				break
			}
			if err != nil {
				break
			}
		}
		strs = strings.Split(strings.Join(strs, ""), ",")
		for _, v := range strs {
			number, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
				break
			}
			out <- number
		}
		close(out)
	}()
	return out, nil
}

func startIntServer() {
	var (
		conn     net.Conn
		err      error
		listener net.Listener
	)
	listener, err = net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err = listener.Accept()
			if err != nil {
				break
			}
			strs := make([]string, 0)
			// 向客户端写出连续的10个整数
			for i := 0; i < 10; i++ {
				strs = append(strs, strconv.Itoa(rand.Intn(100)))
			}
			io.WriteString(conn, strings.Join(strs, ","))
			break
		}

		// 一定要关闭
		conn.Close()
	}()
}
