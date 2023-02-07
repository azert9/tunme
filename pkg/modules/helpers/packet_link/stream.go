package packet_link

import (
	"github.com/azert9/tunme/internal/circular_buffer"
	"github.com/azert9/tunme/pkg/modules"
	"sync"
)

type stream struct {
	id           streamId
	receiver     *streamReceiver
	sendLoop     *streamSendLoop
	packetSender modules.PacketSender
	firstAckOnce sync.Once
	firstAckChan chan struct{} // this chan is closed when the first ACK is received (the connection is established)
}

func newStream(id streamId, packetSender modules.PacketSender) *stream {

	windowLen := 10000 // TODO: configure
	sendingBuff := circular_buffer.NewCircularBuffer(windowLen)
	sendingBuff.SetBlockingWrite(true)
	receivingBuff := circular_buffer.NewCircularBuffer(windowLen)

	s := &stream{
		id:           id,
		sendLoop:     newStreamSendLoop(id, packetSender, sendingBuff),
		packetSender: packetSender,
		firstAckChan: make(chan struct{}),
	}

	s.receiver = newStreamReceiver(id, s, receivingBuff)

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

	s.firstAckOnce.Do(func() {
		close(s.firstAckChan)
	})

	s.sendLoop.handleReceivedAckPacket(packet)
}

func (s *stream) handleReceivedDataPacket(packet dataPacket) {
	s.receiver.handleReceivedDataPacket(packet) // TODO: handle errors
}
