package persistency

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEngine_Compact(t *testing.T) {
	os.Remove("data.txt")
	v1 := []int{-1, 2, 3}
	v2 := []int{-1, 2, -9}
	e, _ := NewEngine("data.txt")
	e.Set("key1", []int{1, 2, 3})
	e.Set("key2", []int{1, 2, 3})
	e.Set("key1", v1)
	e.Set("key2", v1)
	e.Set("key3", v2)

	go e.CompactFile()

	e.Set("key1", []int{1, 2, 3})
	e.Set("key2", []int{1, 2, 3})
	e.Set("key1", v1)
	e.Set("key2", v1)
	e.Set("key3", v2)

	time.Sleep((compactTimer + 3) * time.Second)
	require.Equal(t, e.GetFileContent(e.file), 3)
}

func TestEngine_Restore(t *testing.T) {
	fileName := "data.txt"
	os.Remove(fileName)

	e, err := NewEngine(fileName)
	require.NoError(t, err)

	e.Set("key1_restore", []int{1, 2, 3})
	e.Set("key2_restore", []int{4, 5, 6})

	err = e.Close()
	require.NoError(t, err)

	e, _ = NewEngine(fileName)
	e.Restore()

	value, _ := e.Get("key1_restore")
	require.Equal(t, []int{1, 2, 3}, value)
}
