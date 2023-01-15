package udp

import (
	"fmt"
	"net"
	"sync/atomic"
	"tunme/pkg/helpers/packet_link"
	"tunme/pkg/link"
)

type ServerModule struct {
}

type serverTransport struct {
	conn   *net.UDPConn
	remote atomic.Pointer[net.UDPAddr]
}

func (t *serverTransport) ReceivePacket(out []byte) (int, error) {

	n, addr, err := t.conn.ReadFromUDP(out)
	if err != nil {
		return n, err
	}

	t.remote.Store(addr)

	return n, nil
}

func (t *serverTransport) SendPacket(packet []byte) error {

	_, err := t.conn.WriteToUDP(packet, t.remote.Load())

	return err
}

func (ServerModule) Instantiate(positionalArgs []string, namedArgs []link.NamedArg) (link.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("link options: wrong number of positional arguments")
	}
	serverAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("link options: module does not take any named argument")
	}

	resolvedServerAddress, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", resolvedServerAddress)
	if err != nil {
		return nil, err
	}

	transport := &serverTransport{
		conn: conn,
	}

	return packet_link.NewServer(transport, transport), nil
}
