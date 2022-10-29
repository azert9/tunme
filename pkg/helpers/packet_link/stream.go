package packet_link

import (
	"fmt"
	"tunme/internal/circular_buffer"
)

type stream struct {
	_lastReceivedOffset uint64
	_receivingBuff      *circular_buffer.CircularBuffer
	_ackSender          ackSender
}

func newStream(ackSender ackSender) *stream {
	return &stream{
		_receivingBuff: circular_buffer.NewCircularBuffer(2000), // TODO: configure the capacity
		_ackSender:     ackSender,
	}
}

// handleReceivedPacket should not be called concurrently.
func (s *stream) handleReceivedPacket(packet dataPacket) error {

	if packet.getStreamOffset() > s._lastReceivedOffset {
		return fmt.Errorf("unordered dataPacket")
	}

	if packet.getStreamOffset()+uint64(len(packet.getPayload())) > s._lastReceivedOffset {

		newData := packet.getPayload()[s._lastReceivedOffset-packet.getStreamOffset():]

		if _, err := s._receivingBuff.Write(newData); err != nil {
			return err
		}

		s._lastReceivedOffset += uint64(len(newData))
	}

	if err := s._ackSender.sendAck(0 /*TODO: streamId*/, s._lastReceivedOffset); err != nil {
		return err
	}

	return nil
}

func (s *stream) Read(out []byte) (int, error) {
	return s._receivingBuff.Read(out)
}
