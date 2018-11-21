package common

import (
	"bufio"
	"io"
)

// read file by line
// use callback to do what you want, error will be passed to callback naturally
func ReadLine(rd io.Reader, callback func(line []byte, err error)) {
	bufReader := bufio.NewReader(rd)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		callback(line, err)
	}
}
