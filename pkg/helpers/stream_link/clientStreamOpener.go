package stream_link

import (
	"github.com/azert9/tunme/pkg/link"
	"io"
)

type _clientStreamOpener struct {
	ConnFactory connectionFactory
}

func newClientStreamOpener(connFactory connectionFactory) link.StreamOpener {
	return &_clientStreamOpener{
		ConnFactory: connFactory,
	}
}

func (d *_clientStreamOpener) OpenStream() (io.ReadWriteCloser, error) {
	return d.ConnFactory.MakeConnection(1)
}
