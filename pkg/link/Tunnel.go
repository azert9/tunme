package link

import (
	"io"
)

type PacketTunnel interface {
	SendPacket(packet []byte) error
	ReceivePacket(out []byte) (int, error)
}

type StreamAcceptor interface {
	AcceptStream() (io.ReadWriteCloser, error)
}

type StreamOpener interface {
	OpenStream() (io.ReadWriteCloser, error)
}

type Tunnel interface {
	io.Closer
	PacketTunnel
	StreamAcceptor
	StreamOpener
}
