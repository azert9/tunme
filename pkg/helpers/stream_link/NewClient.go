package stream_link

import (
	"github.com/azert9/tunme/pkg/link"
	"net"
)

type StreamDialer interface {
	Dial() (net.Conn, error)
}

func NewClient(dialer StreamDialer) link.Tunnel {

	connFactory := newClientConnectionFactory(dialer)

	controlSteam := newControlStream(connFactory)

	return &tunnel{
		StreamAcceptor: newClientStreamAcceptor(controlSteam, connFactory),
		StreamOpener:   newClientStreamOpener(connFactory),
		PacketConn:     newPacketConn(controlSteam),
	}
}
