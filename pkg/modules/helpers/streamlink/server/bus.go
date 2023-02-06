package server

import (
	"context"
	"io"
)

type bus struct {
	_closeChan             chan struct{}
	_acceptedStreamsChan   chan io.ReadWriteCloser
	_receivedPacketsChan   chan []byte
	_outControlPacketsChan chan []byte
	_callBackStreamsChan   chan io.ReadWriteCloser
}

func newBus() *bus {
	return &bus{
		_closeChan:             make(chan struct{}),
		_acceptedStreamsChan:   make(chan io.ReadWriteCloser),
		_receivedPacketsChan:   make(chan []byte),
		_outControlPacketsChan: make(chan []byte),
		_callBackStreamsChan:   make(chan io.ReadWriteCloser),
	}
}

func (bus *bus) close() {
	close(bus._closeChan)
}

func (bus *bus) sendAcceptedStreamNonBlocking(stream io.ReadWriteCloser) bool {

	select {
	case bus._acceptedStreamsChan <- stream:
		return true
	case _, _ = <-bus._closeChan:
		return false
	default:
		return false
	}
}

func (bus *bus) receiveAcceptedStream() (io.ReadWriteCloser, bool) {

	select {
	case stream := <-bus._acceptedStreamsChan:
		return stream, true
	case _, _ = <-bus._closeChan:
		return nil, false
	}
}

func (bus *bus) sendReceivedPacket(packet []byte) bool {

	select {
	case bus._receivedPacketsChan <- packet:
		return true
	case _, _ = <-bus._closeChan:
		return false
	}
}

func (bus *bus) receiveReceivedPacket() ([]byte, bool) {

	select {
	case packet := <-bus._receivedPacketsChan:
		return packet, true
	case _, _ = <-bus._closeChan:
		return nil, false
	}
}

func (bus *bus) sendOutControlPacket(packet []byte) bool {

	select {
	case bus._outControlPacketsChan <- packet:
		return true
	case _, _ = <-bus._closeChan:
		return false
	}
}

func (bus *bus) receiveOutControlPacket(ctx context.Context) ([]byte, bool) {

	select {
	case packet := <-bus._outControlPacketsChan:
		return packet, true
	case <-bus._closeChan:
		return nil, false
	case <-ctx.Done():
		return nil, false
	}
}

func (bus *bus) sendCallbackStream(stream io.ReadWriteCloser) bool {

	select {
	case bus._callBackStreamsChan <- stream:
		return true
	case _, _ = <-bus._closeChan:
		return false
	}
}

func (bus *bus) receiveCallbackStream() (io.ReadWriteCloser, bool) {

	select {
	case stream := <-bus._callBackStreamsChan:
		return stream, true
	case _, _ = <-bus._closeChan:
		return nil, false
	}
}
