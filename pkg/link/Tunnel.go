package link

import (
	"io"
)

type PacketSender interface {
	SendPacket(packet []byte) error
}

type PacketReceiver interface {
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
	PacketSender
	PacketReceiver
	StreamAcceptor
	StreamOpener
}
