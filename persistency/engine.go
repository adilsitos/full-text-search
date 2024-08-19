package persistency

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

	if strings.Contains(key, " ") {
		return errors.New("key cannot contain spaces")
	}

	return e.setRaw(key, value)
}

func (e *Engine) setRaw(key string, value []int) error {
	offset, err := e.saveTofile(key, value)
	if err != nil {
		return err
	}

	e.data[key] = offset
	return nil
}

func (e *Engine) saveTofile(key string, value []int) (int64, error) {
	offset, err := e.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	sliceStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(value)), ","), "[]")

	data := fmt.Sprintf("%s%s%d|%s\n", key, keyValueSeprator, len(sliceStr), sliceStr)

	_, err = e.file.WriteString(data)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func (e *Engine) Get(key string) ([]int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	val, ok := e.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	// the + 1 is used to count the space between the key and value on the file
	_, err := e.file.Seek(int64(len(key))+val+1, 0)
	if err != nil {
		return nil, err
	}

	buff := make([]byte, 1)
	var content []byte

	for {
		n, err := e.file.Read(buff)
		if err != nil {
			return nil, err
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
	return DecodeToIntArr(string(content))
}

const compactTimer = 5 // vale in seconds

func (e *Engine) CompactFile() {
	for {
		// todo: change the sleep to a more performant operation
		// try to do it parallel
		time.Sleep(time.Duration(compactTimer) * time.Second)
		fmt.Println("compacting file...")

		e.mu.Lock()

		_, m := e.GetMapFromFile()

		err := e.file.Truncate(0)
		if err != nil {
			fmt.Println(err)
			e.mu.Unlock()
			continue
		}

		for k, v := range m {
			e.setRaw(k, v)
		}

		e.file.Seek(0, 0)
		e.mu.Unlock()
	}
}

type Item struct {
	Key    string
	Value  []int
	Offset int64
}

// it can be better. Try to refactor this to use the same approach of
// lsm trees, where it is possible to have several defined block size
func (e *Engine) GetMapFromFile() ([]Item, map[string][]int) {
	m := make(map[string][]int)
	itens := []Item{}

	_, err := e.file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return itens, m
	}

	var totalBytesRead int64
	scanner := bufio.NewScanner(e.file)

	for scanner.Scan() {
		line := scanner.Text()
		offset := totalBytesRead
		parts := strings.Split(line, keyValueSeprator)
		if len(parts) >= 2 {
			key, value := parts[0], parts[1]

			arr, err := DecodeToIntArr(value)
			if err != nil {
				fmt.Println(err)
				continue
			}

			m[key] = arr
			itens = append(itens, Item{
				Key:    key,
				Value:  arr,
				Offset: offset,
			})
		}
	}

	return itens, m
}

func (e *Engine) GetFileContent(f *os.File) []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := f.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	scanner := bufio.NewScanner(f)
	var content []string

	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	return content
}

func (e *Engine) Restore() {
	e.mu.Lock()
	defer e.mu.Unlock()

	items, _ := e.GetMapFromFile()

	for _, v := range items {
		e.data[v.Key] = v.Offset
	}
}

func (e *Engine) Close() error {
	return e.file.Close()
}

func DecodeToIntArr(val string) ([]int, error) {
	// format = size|val, val, val ...
	if len(val) < 1 {
		return nil, errors.New("invalid string value")
	}

	size, err := strconv.Atoi(string(val[0]))
	if err != nil {
		return nil, err
	}

	// consume | char
	if size != len(val[2:]) {
		return nil, errors.New("invalid value size")
	}

	arrStr := strings.Split(val[2:], ",")

	fmt.Println(arrStr)

	arrInt := make([]int, 0, len(arrStr))
	for _, val := range arrStr {
		aux, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}

		arrInt = append(arrInt, aux)
	}

	return arrInt, nil
}
