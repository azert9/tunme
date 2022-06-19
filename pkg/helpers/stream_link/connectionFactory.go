package stream_link

import (
	"net"
)

type connectionFactory interface {
	MakeConnection(channel int) (net.Conn, error)
}
