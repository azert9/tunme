package stream_link

import (
	"io"
	"tunme/pkg/link"
)

type _serverStreamOpener struct {
	ControlStream controlStream
	ConnFactory   connectionFactory
}

func newServerStreamOpener(controlStream controlStream, connFactory connectionFactory) link.StreamOpener {
	return &_serverStreamOpener{
		ControlStream: controlStream,
		ConnFactory:   connFactory,
	}
}

func (d *_serverStreamOpener) OpenStream() (io.ReadWriteCloser, error) {

	if err := d.ControlStream.SendStreamRequest(); err != nil {
		return nil, err
	}

	return d.ConnFactory.MakeConnection(2)
}
