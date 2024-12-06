package storages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBulkMaxLength(t *testing.T) {
	savedSet := make(map[int]struct{})
	sampleFunc := func(list []interface{}) {
		for _, l := range list {
			savedSet[l.(int)] = struct{}{}
		}
	}
	bulk := newBulk(sampleFunc)

	for i := 0; i < MAX_BULK_LENGTH-1; i++ {
		bulk.append(i)
	}

	assert.Equal(t, 0, len(savedSet))
	bulk.append(999)

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, MAX_BULK_LENGTH, len(savedSet))
}

func TestBulkFlashInterval(t *testing.T) {
	savedSet := make(map[int]struct{})
	sampleFunc := func(list []interface{}) {
		for _, l := range list {
			savedSet[l.(int)] = struct{}{}
		}
	}
	bulk := newBulk(sampleFunc)

	testLength := 30
	assert.Less(t, testLength, MAX_BULK_LENGTH)
	for i := 0; i < testLength; i++ {
		bulk.append(i)
	}

	assert.Equal(t, 0, len(savedSet))

	time.Sleep(AUTO_FLASH_INTERVAL + 1*time.Second)
	assert.Equal(t, testLength, len(savedSet))
}
