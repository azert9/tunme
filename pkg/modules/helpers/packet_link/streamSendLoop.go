package packet_link

import (
	"github.com/azert9/tunme/internal/circular_buffer"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"sync"
	"time"
)

// streamSendLoop is responsible for sending data packets (not ACKs), and handling received ACKs
type streamSendLoop struct {
	buff          *circular_buffer.CircularBuffer
	_streamId     streamId
	_waitGroup    sync.WaitGroup
	_streamOffset uint64 // offset corresponding to the beginning of the buffered data, and the last ack received
	_packetSender modules.PacketSender
	_ackChan      chan uint64
	_closeChan    chan struct{}
}

func newStreamSendLoop(streamId streamId, packetSender modules.PacketSender, buff *circular_buffer.CircularBuffer) *streamSendLoop {

	l := &streamSendLoop{
		buff:          buff,
		_streamId:     streamId,
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

	select {
	case l._ackChan <- packet.getStreamOffset():
	default:
		// ACK received while we were already in the process of retransmitting. Not catastrophic.
	}
}

func (l *streamSendLoop) _loop() {

	defer l._waitGroup.Done()

	// TODO: configure the max payload length according to the max packet size
	packet := newDataPacket(2000)
	packet.setStreamId(l._streamId)

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
