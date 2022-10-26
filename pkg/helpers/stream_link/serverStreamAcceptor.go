package stream_link

import (
	"io"
	"tunme/pkg/link"
)

type _serverStreamAcceptor struct {
	ConnFactory connectionFactory
}

func newServerStreamAcceptor(connFactory connectionFactory) link.StreamAcceptor {
	return &_serverStreamAcceptor{
		ConnFactory: connFactory,
	}
}

func (l *_serverStreamAcceptor) AcceptStream() (io.ReadWriteCloser, error) {
	return l.ConnFactory.MakeConnection(1)
}
