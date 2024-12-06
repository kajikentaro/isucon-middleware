package storages

import (
	"sync"
	"time"
)

var MAX_BULK_LENGTH = 50
var AUTO_FLASH_INTERVAL = 1 * time.Second

type bulk struct {
	chunk     []interface{}
	chunkLock *sync.Mutex
	timer     *time.Timer
	execFunc  func([]interface{})
}

func newBulk(execFunc func([]interface{})) *bulk {
	res := &bulk{
		chunkLock: &sync.Mutex{},
		timer:     time.NewTimer(AUTO_FLASH_INTERVAL),
		execFunc:  execFunc,
	}
	res.startBulkInsertTimer()
	return res
}

func (b *bulk) append(single interface{}) {
	b.chunkLock.Lock()
	defer b.chunkLock.Unlock()

	b.chunk = append(b.chunk, single)

	if len(b.chunk) >= MAX_BULK_LENGTH {
		b.timer.Stop()
		go func() {
			b.execute()
			b.timer.Reset(AUTO_FLASH_INTERVAL)
		}()
	}
}

func (b *bulk) startBulkInsertTimer() {
	go func() {
		for {
			<-b.timer.C
			b.execute()
			b.timer.Reset(AUTO_FLASH_INTERVAL)
		}
	}()
}

func (b *bulk) execute() {
	b.chunkLock.Lock()
	processing := b.chunk
	b.chunk = []interface{}{}
	b.chunkLock.Unlock()

	if len(processing) == 0 {
		return
	}

	b.execFunc(processing)
}
