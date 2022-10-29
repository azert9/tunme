package packet_link

import (
	"tunme/pkg/link"
)

func receiveLoop(receiver link.PacketReceiver, packetsChan chan<- []byte, streams *streamManager) {

	buff := make([]byte, 4096) // TODO: length

	for {

		n, err := receiver.ReceivePacket(buff)
		if err != nil {
			// TODO
			panic(err)
		}

		if n == 0 {
			// keep-alive dataPacket
			continue
		}

		switch buff[0] {
		case byte(packetTypePacket):
			select {
			case packetsChan <- buff[1:n]:
			default:
			}
			buff = make([]byte, len(buff)) // TODO: avoid allocations as much as possible
		case byte(packetTypeStreamData):
			packet, err := dataPacketFromBuff(buff[:n])
			if err != nil {
				panic(err) // TODO
			}
			streams.handleReceivedDataPacket(packet)
		case byte(packetTypeStreamAck):
			packet, err := ackPacketFromBuff(buff[:n])
			if err != nil {
				panic(err) // TODO
			}
			streams.handleReceivedAckPacket(packet)
		}
	}
}
