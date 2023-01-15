package helpers

import (
	"fmt"
	"tunme/pkg/helpers/packet_link"
	"tunme/pkg/helpers/stream_link"
	"tunme/pkg/link"
)

type mockPacketSender chan<- []byte

func (s mockPacketSender) SendPacket(packet []byte) error {
	s <- packet
	return nil
}

type mockPacketReceiver <-chan []byte

func (r mockPacketReceiver) ReceivePacket(out []byte) (int, error) {

	received := <-r

	if len(received) > len(out) {
		return 0, fmt.Errorf("buffer too small")
	}

	copy(out, received)

	return len(received), nil
}

func newMockTunPair(linkType int, roleOrder int) (link.Tunnel, link.Tunnel) {

	var tun1, tun2 link.Tunnel

	if linkType%2 == 0 {

		connectionChan := make(chan mockConn)

		tun1 = stream_link.NewServer(&mockListener{
			incomingConnections: connectionChan,
		})

		tun2 = stream_link.NewClient(&mockDialer{
			outgoingConnections: connectionChan,
		})

	} else {

		// TODO: packetChan2 should become available only after the first message is received on packetChan1
		packetChan1 := make(chan []byte)
		packetChan2 := make(chan []byte)
		tun1 = packet_link.NewClient(mockPacketSender(packetChan1), mockPacketReceiver(packetChan2))
		tun2 = packet_link.NewServer(mockPacketSender(packetChan2), mockPacketReceiver(packetChan1))
	}

	if roleOrder%2 == 0 {
		return tun1, tun2
	} else {
		return tun2, tun1
	}
}
