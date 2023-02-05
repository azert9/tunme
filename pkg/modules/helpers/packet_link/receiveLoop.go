package packet_link

import (
	"github.com/azert9/tunme/pkg/modules"
	"sync"
)

func receiveLoop(receiver modules.PacketReceiver, sender modules.PacketSender, packetsChan chan<- []byte, firstPacketChan chan struct{}, streams *streamManager) {

	var firstPacketOnce sync.Once

	buff := make([]byte, 4096) // TODO: length

	for {

		n, err := receiver.ReceivePacket(buff)
		if err != nil {
			// TODO
			panic(err)
		}

		firstPacketOnce.Do(func() {
			close(firstPacketChan)
		})

		switch buff[0] {
		case byte(packetTypePing):
			var packet [1]byte
			packet[0] = byte(packetTypePong)
			if err := sender.SendPacket(packet[:]); err != nil {
				panic(err) // TODO
			}
		case byte(packetTypePong):
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
