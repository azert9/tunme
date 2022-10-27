package helpers

import (
	"fmt"
	"io"
	"net"
	"time"
)

type mockConn struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (conn mockConn) Read(b []byte) (n int, err error) {
	return conn.reader.Read(b)
}

func (conn mockConn) Write(b []byte) (n int, err error) {
	return conn.writer.Write(b)
}

func (conn mockConn) Close() error {
	conn.reader.Close()
	conn.writer.Close()
	return nil
}

func (conn mockConn) LocalAddr() net.Addr {
	return nil
}

func (conn mockConn) RemoteAddr() net.Addr {
	return nil
}

func (conn mockConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}

func (conn mockConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}

func (conn mockConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("unsupported operation")
}
