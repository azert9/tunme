package stream_link

import (
	"tunme/pkg/helpers/stream_link"
	"tunme/pkg/link"
)

func newMockTunPair(roleOrder int) (link.Tunnel, link.Tunnel) {

	connectionChan := make(chan mockConn)

	tun1 := stream_link.NewServer(&mockListener{
		incomingConnections: connectionChan,
	})

	tun2 := stream_link.NewClient(&mockDialer{
		outgoingConnections: connectionChan,
	})

	if roleOrder%2 == 0 {
		return tun1, tun2
	} else {
		return tun2, tun1
	}
}
