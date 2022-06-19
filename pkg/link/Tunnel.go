package link

import (
	"net"
)

type Listener interface {
	Accept() (net.Conn, error)
}

type StreamDialer interface {
	// TODO: support deadline / timeout
	Dial() (net.Conn, error)
}

// TODO: smoother interface
type Tunnel struct {
	StreamListener net.Listener
	StreamDialer   StreamDialer
	PacketConn     net.PacketConn
}
