package common

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func IsDir(name string) bool {
	f, err := os.Stat(name)
	if err != nil {
		return false
	}

	return f.IsDir()
}

// get curr dir of running go program
func GetCWD() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func MustGetCWD() string {
	cwd, err := GetCWD()
	if err != nil {
		panic(err)
	}
	return cwd
}

// a.txt {"txt", ".jpg"} // in targets, dot(.) is optional
func HasSuffix(filename string, targets []string) bool {

	// 文件名中没有 .  认为是false
	if strings.IndexAny(filename, ".") < 0 {
		return false
	}

	for _, suffix := range targets {
		tmp := suffix
		if tmp[0] != '.' { // such as "txt" instead of ".txt"
			tmp = "." + tmp
		}

		if filepath.Ext(filename) == tmp {
			return true
		}
	}
	return false
}
