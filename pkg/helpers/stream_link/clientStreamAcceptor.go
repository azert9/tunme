package stream_link

import (
	"github.com/azert9/tunme/pkg/link"
	"io"
)

type _clientStreamAcceptor struct {
	ControlStream controlStream
	ConnFactory   connectionFactory
}

func newClientStreamAcceptor(controlStream controlStream, connFactory connectionFactory) link.StreamAcceptor {
	return &_clientStreamAcceptor{
		ControlStream: controlStream,
		ConnFactory:   connFactory,
	}
}

func (l *_clientStreamAcceptor) AcceptStream() (io.ReadWriteCloser, error) {

	l.ControlStream.WaitStreamRequest()

	return l.ConnFactory.MakeConnection(2)
}
