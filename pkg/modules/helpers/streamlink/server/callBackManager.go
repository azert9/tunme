package server

import (
	"fmt"
	"io"
	"sync"
)

// TODO: remove

type callBackManager struct {
	mutex    sync.Mutex
	channels map[uint64]chan io.ReadWriteCloser
	nextId   uint64
}

func (cbm *callBackManager) subscribe() (uint64, <-chan io.ReadWriteCloser) {

	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()

	id := cbm.nextId
	cbm.nextId++

	c := make(chan io.ReadWriteCloser, 1)

	cbm.channels[id] = c

	return id, c
}

func (cbm *callBackManager) giveStream(id uint64, stream io.ReadWriteCloser) error {

	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()

	c, ok := cbm.channels[id]
	if !ok {
		return fmt.Errorf("no such stream request")
	}

	delete(cbm.channels, id)

	c <- stream

	return nil
}
