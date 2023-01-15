package packet_link

import (
	"github.com/azert9/tunme/pkg/link"
)

func NewClient(sender link.PacketSender, receiver link.PacketReceiver) link.Tunnel {
	return newTunnel(sender, receiver, false)
}
