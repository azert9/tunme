package udp

import (
	"fmt"
	"net"
	"tunme/pkg/helpers/packet_link"
	"tunme/pkg/link"
)

type ClientModule struct {
}

type clientTransport struct {
	conn *net.UDPConn
}

func (t *clientTransport) ReceivePacket(out []byte) (int, error) {
	return t.conn.Read(out)
}

func (t *clientTransport) SendPacket(packet []byte) error {
	_, err := t.conn.Write(packet)
	return err
}

func (ClientModule) Instantiate(positionalArgs []string, namedArgs []link.NamedArg) (link.Tunnel, error) {

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

	conn, err := net.DialUDP("udp", nil, resolvedServerAddress)
	if err != nil {
		return nil, err
	}

	transport := &clientTransport{
		conn: conn,
	}

	return packet_link.NewClient(transport, transport), nil
}
