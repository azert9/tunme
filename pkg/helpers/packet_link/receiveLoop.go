package packet_link

import (
	"log"
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
			// keep-alive packet
			continue
		}

		switch buff[0] {
		case 0: // packet
			select {
			case packetsChan <- buff[1:n]:
			default:
			}
			buff = make([]byte, len(buff)) // TODO: avoid allocations as much as possible
		case 1: // stream
			streams.HandlePacket(buff[:n])
		}

		if buff[0] == 0 {
		} else {
			// TODO
			log.Printf("invalid packet type")
		}
	}
}
