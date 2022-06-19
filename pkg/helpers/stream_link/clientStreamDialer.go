package stream_link

import (
	"net"
	"tunme/pkg/link"
)

type _clientStreamDialer struct {
	ConnFactory connectionFactory
}

func newClientStreamDialer(connFactory connectionFactory) link.StreamDialer {
	return &_clientStreamDialer{
		ConnFactory: connFactory,
	}
}

func (d *_clientStreamDialer) Dial() (net.Conn, error) {
	return d.ConnFactory.MakeConnection(1)
}
