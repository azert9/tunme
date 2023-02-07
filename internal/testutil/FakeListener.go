package testutil

import "net"

type FakeListener struct {
	IncomingConnections <-chan FakeConn
}

func (l *FakeListener) Accept() (net.Conn, error) {
	conn := <-l.IncomingConnections
	return conn, nil
}

func (l *FakeListener) Close() error {
	//TODO implement me
	panic("implement me")
}

func (l *FakeListener) Addr() net.Addr {
	return nil
}
