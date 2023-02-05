package testutil

import (
	"io"
	"net"
)

type FakeDialer struct {
	OutgoingConnections chan<- FakeConn
}

func (d *FakeDialer) Dial() (net.Conn, error) {

	pipe1r, pipe1w := io.Pipe()
	pipe2r, pipe2w := io.Pipe()

	ownConn := FakeConn{
		reader: pipe1r,
		writer: pipe2w,
	}

	otherConn := FakeConn{
		reader: pipe2r,
		writer: pipe1w,
	}

	d.OutgoingConnections <- otherConn

	return ownConn, nil
}
