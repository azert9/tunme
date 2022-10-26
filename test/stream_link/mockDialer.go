package stream_link

import (
	"io"
	"net"
)

type mockDialer struct {
	outgoingConnections chan<- mockConn
}

func (d *mockDialer) Dial() (net.Conn, error) {

	pipe1r, pipe1w := io.Pipe()
	pipe2r, pipe2w := io.Pipe()

	ownConn := mockConn{
		reader: pipe1r,
		writer: pipe2w,
	}

	otherConn := mockConn{
		reader: pipe2r,
		writer: pipe1w,
	}

	d.outgoingConnections <- otherConn

	return ownConn, nil
}
