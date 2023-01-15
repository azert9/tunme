package packet_link

import "github.com/azert9/tunme/pkg/link"

func NewServer(sender link.PacketSender, receiver link.PacketReceiver) link.Tunnel {
	return newTunnel(sender, receiver, true)
}
