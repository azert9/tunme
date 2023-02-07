package udp

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/modules/helpers/packet_link"
	"net"
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

func (ClientModule) Instantiate(positionalArgs []string, namedArgs []modules.NamedArg) (modules.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("modules options: wrong number of positional arguments")
	}
	serverAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("modules options: modules does not take any named argument")
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
