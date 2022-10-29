package packet_link

import (
	"tunme/internal/circular_buffer"
	"tunme/pkg/link"
)

type stream struct {
	receiver     *streamReceiver
	sendLoop     *streamSendLoop
	packetSender link.PacketSender
}

func newStream(packetSender link.PacketSender) *stream {

	windowLen := 65536 // TODO: configure
	sendingBuff := circular_buffer.NewCircularBuffer(windowLen)
	receivingBuff := circular_buffer.NewCircularBuffer(windowLen)

	s := &stream{
		sendLoop:     newStreamSendLoop(packetSender, sendingBuff),
		packetSender: packetSender,
	}

	s.receiver = newStreamReceiver(s, receivingBuff)

	return s
}

func (s *stream) Close() error {

	// TODO: properly close the stream on both ends

	s.sendLoop.close() // TODO: do not call twice

	return nil
}

func (s *stream) Read(out []byte) (int, error) {
	return s.receiver.buff.Read(out)
}

func (s *stream) Write(buff []byte) (int, error) {
	return s.sendLoop.buff.Write(buff)
}

func (s *stream) sendAck(streamId streamId, offset uint64) error {

	packet := newAckPacket()
	packet.setStreamId(streamId)
	packet.setStreamOffset(offset)

	return s.packetSender.SendPacket(packet.getBytes())
}

func (s *stream) handleReceivedAckPacket(packet ackPacket) {
	s.sendLoop.handleReceivedAckPacket(packet)
}

func (s *stream) handleReceivedDataPacket(packet dataPacket) {
	s.receiver.handleReceivedDataPacket(packet) // TODO: handle errors
}
