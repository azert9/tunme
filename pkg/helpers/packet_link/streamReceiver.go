package packet_link

import (
	"fmt"
	"tunme/internal/circular_buffer"
)

// streamReceiver is responsible for handling incoming data packets and sending ACKs
type streamReceiver struct {
	buff                *circular_buffer.CircularBuffer
	_lastReceivedOffset uint64
	_ackSender          ackSender
}

func newStreamReceiver(ackSender ackSender, buff *circular_buffer.CircularBuffer) *streamReceiver {
	return &streamReceiver{
		buff:       buff,
		_ackSender: ackSender,
	}
}

// handleReceivedDataPacket should not be called concurrently.
func (s *streamReceiver) handleReceivedDataPacket(packet dataPacket) error {

	if packet.getStreamOffset() > s._lastReceivedOffset {
		return fmt.Errorf("unordered dataPacket")
	}

	if packet.getStreamOffset()+uint64(len(packet.getPayload())) > s._lastReceivedOffset {

		newData := packet.getPayload()[s._lastReceivedOffset-packet.getStreamOffset():]

		if _, err := s.buff.Write(newData); err != nil {
			return err
		}

		s._lastReceivedOffset += uint64(len(newData))
	}

	if err := s._ackSender.sendAck(0 /*TODO: streamId*/, s._lastReceivedOffset); err != nil {
		return err
	}

	return nil
}
