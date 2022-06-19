package stream_link

import (
	"net"
	"tunme/pkg/link"
)

func NewServer(listener net.Listener) *link.Tunnel {

	connFactory := newServerConnectionFactory(listener, 3)

	controlSteam := newControlStream(connFactory)

	return &link.Tunnel{
		StreamListener: newServerStreamListener(connFactory),
		StreamDialer:   newServerStreamDialer(controlSteam, connFactory),
		PacketConn:     newPacketConn(controlSteam),
	}
}
