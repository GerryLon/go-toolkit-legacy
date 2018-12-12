package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

// multi goroutine copy
// 多个goroutine复制文件
// cp -n 2 -s src -d dst
func main() {
	n := flag.Int("n", runtime.NumCPU(), "number of goroutine, default is: runtime.NumCPU()")
	s := flag.String("s", "", "source file")
	d := flag.String("d", "", "destination file")

	flag.Parse()

	// check arguments
	if *n < 1 {
		*n = 1
	}

	if len(strings.TrimSpace(*s)) == 0 {
		fmt.Printf("source file can not be empty\n")
		os.Exit(1)
	}

	if len(strings.TrimSpace(*d)) == 0 {
		fmt.Printf("destination file can not be empty\n")
		os.Exit(1)
	}

	sfd, err := os.Open(*s)
	if err != nil {
		fmt.Printf("open source file, err: %v\n", err)
		os.Exit(1)
	}
	defer sfd.Close()
	info, err := sfd.Stat()

	if err != nil {
		fmt.Printf("stat source file, err: %v\n", err)
		os.Exit(1)
	}

	// length in bytes
	size := info.Size()
	var block int64 = size / int64(*n)
	extra := size % block

	// 0644 TODO: dst file mode should be dynamic
	flags := os.O_CREATE | os.O_EXCL | os.O_WRONLY

	dfd, err := os.OpenFile(*d, flags, 0644)

	if err != nil {
		fmt.Printf("open destination file, err: %v\n", err)
		os.Exit(1)
	}
	defer dfd.Close()

	// copy file with multi goroutine
	var wg sync.WaitGroup
	wg.Add(*n)
	for i := 0; i < *n; i++ {
		go func(i int64) {
			var extra2 int64 = 0
			if i == int64(*n-1) {
				extra2 = extra
			}
			cp(sfd, dfd, i*block, extra2, block)
			wg.Done()
		}(int64(i))
	}
	wg.Wait()
}

func cp(sfd, dfd *os.File, offset, extra, block int64) {

	buf := make([]byte, block+extra)
	for {
		n, err := sfd.ReadAt(buf, offset)

		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Printf("read from source file, err: %v\n", err)
				os.Exit(1)
			}
		} else if n == 0 { // finished
			return
		}

		n, err = dfd.WriteAt(buf, offset)
		if err != nil {
			fmt.Printf("write to destination file, err: %v\n", err)
			os.Exit(1)
		}
		return
	}
}
