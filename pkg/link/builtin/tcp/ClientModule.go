package tcp

import (
	"fmt"
	"net"
	"tunme/pkg/helpers/stream_link"
	"tunme/pkg/link"
)

type ClientModule struct {
}

type _streamDialer struct {
	ServerAddress string
}

func (d *_streamDialer) Dial() (net.Conn, error) {
	return net.Dial("tcp", d.ServerAddress)
}

func (ClientModule) Instantiate(positionalArgs []string, namedArgs []link.NamedArg) (*link.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("link options: wrong number of positional arguments")
	}
	serverAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("link options: module does not take any namd argument")
	}

	return stream_link.NewClient(&_streamDialer{serverAddress}), nil
}
