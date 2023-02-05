package tcp

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/modules/helpers/streamlink/server"
	"net"
)

type ServerModule struct {
}

func (ServerModule) Instantiate(positionalArgs []string, namedArgs []modules.NamedArg) (modules.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("modules options: wrong number of positional arguments")
	}
	listenAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("modules options: modules does not take any namd argument")
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	return server.NewTunnel(listener), nil
}
