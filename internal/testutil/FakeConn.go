package testutil

import (
	"fmt"
	"io"
	"net"
	"time"
)

type FakeConn struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (conn FakeConn) Read(b []byte) (n int, err error) {
	return conn.reader.Read(b)
}

func (conn FakeConn) Write(b []byte) (n int, err error) {
	return conn.writer.Write(b)
}

func (conn FakeConn) Close() error {
	conn.reader.Close()
	conn.writer.Close()
	return nil
}

func (conn FakeConn) LocalAddr() net.Addr {
	return nil
}

func (conn FakeConn) RemoteAddr() net.Addr {
	return nil
}

func (conn FakeConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}

func (conn FakeConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}

func (conn FakeConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}
