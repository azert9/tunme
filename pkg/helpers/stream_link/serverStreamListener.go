package stream_link

import (
	"net"
)

type _serverStreamListener struct {
	ConnFactory connectionFactory
}

func newServerStreamListener(connFactory connectionFactory) net.Listener {
	return &_serverStreamListener{
		ConnFactory: connFactory,
	}
}

func (l *_serverStreamListener) Accept() (net.Conn, error) {
	return l.ConnFactory.MakeConnection(1)
}

func (l *_serverStreamListener) Addr() net.Addr {
	// TODO
	panic("unimplemented")
}

func (l *_serverStreamListener) Close() error {
	// TODO
	panic("unimplemented")
}
