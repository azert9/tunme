package packet_link

import (
	"fmt"
	"io"
	"time"
	"tunme/pkg/link"
)

type tunnel struct {
	sender          link.PacketSender
	packetChan      chan []byte
	streams         *streamManager
	firstPacketChan chan struct{} // closed when the first packet is received
}

func newTunnel(sender link.PacketSender, receiver link.PacketReceiver, isServer bool) link.Tunnel {

	// TODO: the client should send regular keep-alive packets

	tun := &tunnel{
		sender:          sender,
		packetChan:      make(chan []byte, 16), // TODO: configure capacity
		streams:         newStreamManager(isServer, sender),
		firstPacketChan: make(chan struct{}),
	}

	go receiveLoop(receiver, sender, tun.packetChan, tun.firstPacketChan, tun.streams)

	if isServer {
		_ = <-tun.firstPacketChan
	} else {

		var packet [1]byte
		packet[0] = byte(packetTypePing)

		done := false
		for !done {

			if err := sender.SendPacket(packet[:]); err != nil {
				panic(err) // TODO
			}

			select {
			case _ = <-tun.firstPacketChan:
				done = true
			case _ = <-time.After(4 * time.Second): // TODO: configure
			}
		}
	}

	return tun
}

func (tun *tunnel) Close() error {
	//TODO implement me
	panic("implement me")
}

func (tun *tunnel) SendPacket(packet []byte) error {

	crafted := make([]byte, 1+len(packet))
	copy(crafted[1:], packet)

	return tun.sender.SendPacket(crafted)
}

func (tun *tunnel) ReceivePacket(out []byte) (int, error) {

	packet := <-tun.packetChan

	if len(out) < len(packet) {
		return 0, fmt.Errorf("dataPacket too large for buffer")
	}

	copy(out, packet)

	return len(packet), nil
}

func (tun *tunnel) AcceptStream() (io.ReadWriteCloser, error) {
	return tun.streams.accept()
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {
	return tun.streams.openStream(), nil
}
