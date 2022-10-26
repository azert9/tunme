package stream_link

import (
	"io"
	"net"
)

type tunnel struct {
	PacketConn     net.PacketConn
	StreamListener net.Listener
	StreamDialer   StreamDialer
}

func (tun *tunnel) Close() error {
	//TODO implement me
	panic("implement me")
}

func (tun *tunnel) SendPacket(packet []byte) error {
	_, err := tun.PacketConn.WriteTo(packet, nil)
	return err
}

func (tun *tunnel) ReceivePacket(out []byte) (int, error) {
	n, _, err := tun.PacketConn.ReadFrom(out)
	return n, err
}

func (tun *tunnel) AcceptStream() (io.ReadWriteCloser, error) {
	return tun.StreamListener.Accept()
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {
	return tun.StreamDialer.Dial()
}
