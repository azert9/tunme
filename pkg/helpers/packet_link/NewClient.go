package packet_link

import "tunme/pkg/link"

func NewServer(sender link.PacketSender, receiver link.PacketReceiver) link.Tunnel {
	return newTunnel(sender, receiver)
}
