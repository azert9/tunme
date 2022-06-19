package stream_link

import (
	"net"
	"tunme/pkg/link"
)

type _serverStreamDialer struct {
	ControlStream controlStream
	ConnFactory   connectionFactory
}

func newServerStreamDialer(controlStream controlStream, connFactory connectionFactory) link.StreamDialer {
	return &_serverStreamDialer{
		ControlStream: controlStream,
		ConnFactory:   connFactory,
	}
}

func (d *_serverStreamDialer) Dial() (net.Conn, error) {

	if err := d.ControlStream.SendStreamRequest(); err != nil {
		return nil, err
	}

	return d.ConnFactory.MakeConnection(2)
}
