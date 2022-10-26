package link

import (
	"io"
	"net"
)

type Listener interface {
	Accept() (net.Conn, error)
}

type StreamDialer interface {
	// TODO: support deadline / timeout
	Dial() (net.Conn, error)
}

type PacketTunnel interface {
	SendPacket(packet []byte) error
	ReceivePacket(out []byte) (int, error)
}

type StreamTunnel interface {
	AcceptStream() (io.ReadWriteCloser, error)
	OpenStream() (io.ReadWriteCloser, error)
}

type Tunnel interface {
	io.Closer
	PacketTunnel
	StreamTunnel
}
