package persistency

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// this package is responsible to save data on disk

type Engine struct {
	data map[string]int64 // the value stores the offset location of the file
	file *os.File
	mu   sync.Mutex
}

var keyValueSeprator = " "

func NewEngine(fileName string) (*Engine, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Engine{
		file: file,
		mu:   sync.Mutex{},
		data: make(map[string]int64),
	}, nil
}

// set will be an append only
func (e *Engine) Set(key string, value []int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// set the offset to the end of the file
	offset, err := e.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	sliceStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(value)), ","), "[]")

	data := fmt.Sprintf("%s%s%d%s", key, keyValueSeprator, len(value), sliceStr)

	_, err = e.file.WriteString(data)
	if err != nil {
		return err
	}

	e.data[key] = offset
	return nil
}

func (e *Engine) Get(key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, ok := e.data[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	_, err := e.file.Seek(int64(len(key))+val+1, 0)
	if err != nil {
		return "", err
	}

	buff := make([]byte, 1)
	var content []byte

	for {
		n, err := e.file.Read(buff)
		if err != nil {
			return "", err
		}

		// returns 0 if it's EOF or if did not read any byte
		if n == 0 {
			break
		}

		if buff[0] == '\n' {
			break
		}

		content = append(content, buff[0])
	}

	// TODO: decode the slice of ints (encoded as string) to slice of ints []int
	return string(content), nil
}
