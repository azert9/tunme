package packet_link

type ackSender interface {
	sendAck(streamId streamId, offset uint64) error
}
