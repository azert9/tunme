package helpers

import "net"

type mockListener struct {
	incomingConnections <-chan mockConn
}

func (l *mockListener) Accept() (net.Conn, error) {
	conn := <-l.incomingConnections
	return conn, nil
}

func (l *mockListener) Close() error {
	//TODO implement me
	panic("implement me")
}

func (l *mockListener) Addr() net.Addr {
	return nil
}
