package client

import (
	"io"
)

// TODO: use something similar for the server

type bus struct {
	_closeChan  chan struct{}
	_streamChan chan io.ReadWriteCloser
}

func newBus() *bus {
	return &bus{
		_closeChan:  make(chan struct{}),
		_streamChan: make(chan io.ReadWriteCloser),
	}
}

func (bus *bus) close() {
	close(bus._closeChan)
}

func (bus *bus) sendStream(stream io.ReadWriteCloser) bool {

	select {
	case bus._streamChan <- stream:
		return true
	case _, _ = <-bus._closeChan:
		return false
	}
}

func (bus *bus) receiveStream() (io.ReadWriteCloser, bool) {

	select {
	case stream := <-bus._streamChan:
		return stream, true
	case _, _ = <-bus._closeChan:
		return nil, false
	}
}
