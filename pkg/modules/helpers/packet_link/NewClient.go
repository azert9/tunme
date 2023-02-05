package packet_link

import (
	"github.com/azert9/tunme/pkg/modules"
)

func NewClient(sender modules.PacketSender, receiver modules.PacketReceiver) modules.Tunnel {
	return newTunnel(sender, receiver, false)
}
