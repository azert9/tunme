package stream_link

import "net"

type _clientStreamListener struct {
	ControlStream controlStream
	ConnFactory   connectionFactory
}

func newClientStreamListener(controlStream controlStream, connFactory connectionFactory) net.Listener {
	return &_clientStreamListener{
		ControlStream: controlStream,
		ConnFactory:   connFactory,
	}
}

func (l *_clientStreamListener) Accept() (net.Conn, error) {

	l.ControlStream.WaitStreamRequest()

	return l.ConnFactory.MakeConnection(2)
}

func (*_clientStreamListener) Addr() net.Addr {
	// TODO
	panic("unimplemented")
}

func (*_clientStreamListener) Close() error {
	// TODO
	panic("unimplemented")
}
