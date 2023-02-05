package tcp

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/modules/helpers/streamlink/client"
	"io"
	"net"
)

type ClientModule struct {
}

type _streamDialer struct {
	ServerAddress string
}

func (d *_streamDialer) Dial() (io.ReadWriteCloser, error) {
	return net.Dial("tcp", d.ServerAddress)
}

func (ClientModule) Instantiate(positionalArgs []string, namedArgs []modules.NamedArg) (modules.Tunnel, error) {

	if len(positionalArgs) != 1 {
		return nil, fmt.Errorf("modules options: wrong number of positional arguments")
	}
	serverAddress := positionalArgs[0]

	if len(namedArgs) != 0 {
		return nil, fmt.Errorf("modules options: modules does not take any namd argument")
	}

	return client.NewTunnel(&_streamDialer{serverAddress}), nil
}
