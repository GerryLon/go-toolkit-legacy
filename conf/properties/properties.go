package properties

import (
	"errors"
	"fmt"
	"github.com/GerryLon/go-toolkit/common"
	"log"
	"os"
	"strings"
	"sync"
)

// .properties文件读写
type Properties struct {
	Filename   string // .properties file name
	properties map[string]string
	loaded     bool
	lock       sync.RWMutex
}

func init() {

}

// get value by given key
func (p *Properties) Value(key string) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.loaded {
		err := p.load()
		if err != nil {
			return "", err
		}
	}

	val, ok := p.properties[key]
	if !ok {
		return "", fmt.Errorf("key %s is not exiest in %s", key, p.Filename)
	}
	return val, nil
}

// get value by given key, if error occurred, panic
func (p *Properties) MustValue(key string) (val string) {
	val, err := p.Value(key)
	if err != nil {
		panic(err)
	}
	return val
}

// get all
func (p *Properties) All() (map[string]string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.loaded {
		err := p.load()
		if err != nil {
			return nil, err
		}
	}
	return p.properties, nil
}

// get all, if error occurred, panic
func (p *Properties) MustAll() map[string]string {
	all, err := p.All()
	if err != nil {
		panic(err)
	}
	return all
}

func (p *Properties) load() error {
	if p.loaded {
		return fmt.Errorf("file: %s is already loaded", p.Filename)
	}

	file, err := os.Open(p.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	p.properties = make(map[string]string)
	common.ReadLine(file, func(line []byte, err error) {
		if err != nil {
			log.Printf("common.ReadLine, err: %v", err)
			return
		}

		key, val, err := parse(string(line))
		if err != nil {
			log.Printf("Properties parse(), err: %v", err)
			return
		}

		p.properties[key] = val
	})

	p.loaded = true
	return nil
}

func parse(line string) (key, val string, err error) {
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		return key, val, errors.New("empty line")
	}

	if strings.HasPrefix(line, "#") {
		return key, val, errors.New("comment line")
	}

	index := strings.Index(line, "=")
	if index < 0 {
		return key, val, fmt.Errorf("the line %s dosen't contains =", line)
	}

	key = strings.TrimSpace(line[:index])
	val = strings.TrimSpace(line[index+1:])
	err = nil
	return
}
