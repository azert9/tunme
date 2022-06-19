package stream_link

import (
	"errors"
	"io"
	"net"
	"time"
)

// TODO: could use multiple parallel connections
// TODO: could avoid waiting for ACK
// TODO: errors should close the stream, or trigger a re-open

var errPacketTooBig = errors.New("packet is too big")

type _packetConn struct {
	controlStream controlStream
}

func newPacketConn(controlStream controlStream) net.PacketConn {
	return &_packetConn{
		controlStream: controlStream,
	}
}

func (c *_packetConn) ReadFrom(p []byte) (n int, addr net.Addr, readErr error) {

	packet, readErr := c.controlStream.Receive()

	truncated := false

	copy(p, packet)
	n = len(packet)
	if len(p) < n {
		truncated = true
		n = len(p)
	}

	if readErr == nil && truncated {
		readErr = io.ErrShortBuffer
	}

	return n, nil, readErr
}

func (c *_packetConn) WriteTo(p []byte, _ net.Addr) (n int, err error) {

	if err := c.controlStream.Send(p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (c *_packetConn) Close() error {
	//TODO implement me
	panic("implement me")
}

func (c *_packetConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c *_packetConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *_packetConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *_packetConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}
