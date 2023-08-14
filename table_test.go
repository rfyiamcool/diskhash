package diskhash

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func destroyTable(t *Table) {
	_ = t.Close()
	_ = os.RemoveAll(t.options.DirPath)
}

func GetTestKey(i int) []byte {
	return []byte(fmt.Sprintf("diskhash-test-key-%09d", i))
}

func TestOpen(t *testing.T) {
	dir, err := os.MkdirTemp("", "diskhash-test-open")
	assert.Nil(t, err)

	options := DefaultOptions
	options.DirPath = dir
	options.SlotValueLength = 10
	table, err := Open(options)
	assert.Nil(t, err)
	defer destroyTable(table)

	err = table.Close()
	assert.Nil(t, err)
}

func TestTable_Put(t *testing.T) {
	t.Run("16B-10000000", func(t *testing.T) {
		testTablePut(t, 16, 10000000)
	})
	t.Run("20B-2000000", func(t *testing.T) {
		testTablePut(t, 16, 2000000)
	})
	t.Run("1K-500000", func(t *testing.T) {
		testTablePut(t, 1024, 500000)
	})
	t.Run("4K-500000", func(t *testing.T) {
		testTablePut(t, 4*1024, 500000)
	})
}

func testTablePut(t *testing.T, valueLen uint32, count int) {
	dir, err := os.MkdirTemp("", "diskhash-test-put")
	assert.Nil(t, err)

	options := DefaultOptions
	options.DirPath = dir
	options.SlotValueLength = valueLen
	table, err := Open(options)
	assert.Nil(t, err)
	defer destroyTable(table)

	for i := 0; i < count; i++ {
		key := GetTestKey(i)
		value := []byte(strings.Repeat("D", int(valueLen)))

		err = table.Put(key, value, func(slot Slot) (bool, error) {
			return false, nil
		})

		assert.Nil(t, err)
	}
}

func TestTableCrud(t *testing.T) {
	dir, err := os.MkdirTemp("", "diskhash-test-crud")
	assert.Nil(t, err)

	options := DefaultOptions
	options.DirPath = dir
	options.SlotValueLength = 32
	table, err := Open(options)
	assert.Nil(t, err)
	defer destroyTable(table)

	for i := 0; i < 100; i++ {
		var cur []byte

		getFunc := func(slot Slot) (bool, error) {
			cur = slot.Value
			return false, nil
		}
		updateFunc := func(slot Slot) (bool, error) {
			return false, nil
		}

		key := GetTestKey(i)
		value := []byte(strings.Repeat("D", 32))

		// put
		err = table.Put(key, value, updateFunc)
		assert.Nil(t, err)

		// get
		err = table.Get(key, getFunc)
		assert.Nil(t, err)
		assert.Equal(t, value, cur)

		// put different value
		value = []byte(strings.Repeat("A", 32))
		err = table.Put(key, value, updateFunc)
		assert.Nil(t, err)

		// get after put different value
		err = table.Get(key, getFunc)
		assert.Nil(t, err)
		assert.Equal(t, value, cur)

		// delete
		err = table.Delete(key, updateFunc)
		assert.Nil(t, err)
	}
}
