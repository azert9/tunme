package stream_link

import (
	"io"
	"net"
	"tunme/pkg/link"
)

type tunnel struct {
	PacketConn     net.PacketConn
	StreamAcceptor link.StreamAcceptor
	StreamOpener   link.StreamOpener
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
	return tun.StreamAcceptor.AcceptStream()
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {
	return tun.StreamOpener.OpenStream()
}
