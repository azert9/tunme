package packet_link

import (
	"io"
	"sync"
	"time"
	"tunme/internal/circular_buffer"
	"tunme/pkg/link"
)

// streamSendLoop is responsible for sending data packets (not ACKs), and handling received ACKs
type streamSendLoop struct {
	buff          *circular_buffer.CircularBuffer
	_waitGroup    sync.WaitGroup
	_streamOffset uint64 // offset corresponding to the beginning of the buffered data, and the last ack received
	_packetSender link.PacketSender
	_ackChan      chan uint64
	_closeChan    chan struct{}
}

func newStreamSendLoop(packetSender link.PacketSender, buff *circular_buffer.CircularBuffer) *streamSendLoop {

	l := &streamSendLoop{
		buff:          buff,
		_packetSender: packetSender,
		_ackChan:      make(chan uint64),
		_closeChan:    make(chan struct{}),
	}

	l._waitGroup.Add(1)
	go l._loop()

	return l
}

func (l *streamSendLoop) close() {

	// TODO: ensure not called twice

	_ = l.buff.Close()

	close(l._closeChan)

	l._waitGroup.Wait()
}

func (l *streamSendLoop) handleReceivedAckPacket(packet ackPacket) {
	l._ackChan <- packet.getStreamOffset()
}

func (l *streamSendLoop) _loop() {

	defer l._waitGroup.Done()

	// TODO: configure the max payload length according to the max packet size
	packet := newDataPacket(2000)

	// TODO: set streamId

	ackTimeout := 4 * time.Second // TODO: configure

	for {

		n, err := l.buff.Peek(packet.getPayload())
		if err == io.EOF {
			if n == 0 {
				break
			}
		} else if err != nil {
			panic(err) // unexpected error
		}

		packet.setStreamOffset(l._streamOffset)

		if err := l._packetSender.SendPacket(packet.getBytes(n)); err != nil {
			// TODO
			panic(err)
		}

		select {
		case ackOffset := <-l._ackChan:
			if ackOffset > l._streamOffset {
				transmittedLen := ackOffset - l._streamOffset
				if transmittedLen > uint64(l.buff.Len()) {
					panic("ack for data that is not even available in the sending buffer") // TODO
				}
				_, err := io.ReadFull(l.buff, packet.getPayload()[:transmittedLen])
				if err != nil {
					panic(err) // unexpected error
				}
				l._streamOffset += transmittedLen
			}
		case _ = <-l._closeChan:
			return
		case _ = <-time.After(ackTimeout):
		}
	}
}
