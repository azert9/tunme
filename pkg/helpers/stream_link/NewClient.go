package stream_link

import (
	"net"
	"tunme/pkg/link"
)

type StreamDialer interface {
	Dial() (net.Conn, error)
}

func NewClient(dialer StreamDialer) link.Tunnel {

	connFactory := newClientConnectionFactory(dialer)

	controlSteam := newControlStream(connFactory)

	return &tunnel{
		StreamListener: newClientStreamListener(controlSteam, connFactory),
		StreamDialer:   newClientStreamDialer(connFactory),
		PacketConn:     newPacketConn(controlSteam),
	}
}
