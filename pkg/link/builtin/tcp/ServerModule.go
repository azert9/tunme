package tcp

import (
	"fmt"
	"github.com/azert9/tunme/pkg/helpers/stream_link"
	"github.com/azert9/tunme/pkg/link"
	"net"
)

type ServerModule struct {
}

func (ServerModule) Instantiate(positionalArgs []string, namedArgs []link.NamedArg) (link.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("link options: wrong number of positional arguments")
	}
	listenAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("link options: module does not take any namd argument")
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	return stream_link.NewServer(listener), nil
}
