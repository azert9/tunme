package stream_link

import (
	"github.com/azert9/tunme/pkg/link"
	"net"
)

func NewServer(listener net.Listener) link.Tunnel {

	connFactory := newServerConnectionFactory(listener, 3)

	controlSteam := newControlStream(connFactory)

	return &tunnel{
		StreamAcceptor: newServerStreamAcceptor(connFactory),
		StreamOpener:   newServerStreamOpener(controlSteam, connFactory),
		PacketConn:     newPacketConn(controlSteam),
	}
}
