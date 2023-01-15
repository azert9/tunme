package stream_link

import (
	"github.com/azert9/tunme/pkg/link"
	"io"
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
