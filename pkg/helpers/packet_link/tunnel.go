package packet_link

import (
	"fmt"
	"io"
	"tunme/pkg/link"
)

type tunnel struct {
	sender     link.PacketSender
	packetChan chan []byte
	streams    *streamManager
}

func newTunnel(sender link.PacketSender, receiver link.PacketReceiver, isServer bool) link.Tunnel {

	// TODO: the client should send a first packet to traverse NATs

	tun := &tunnel{
		sender:     sender,
		packetChan: make(chan []byte, 16), // TODO: configure capacity
		streams:    newStreamManager(isServer, sender),
	}

	go receiveLoop(receiver, tun.packetChan, tun.streams)

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
