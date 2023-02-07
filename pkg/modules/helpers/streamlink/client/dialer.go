package client

import (
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"io"
)

type dialerWrapper struct {
	dialer Dialer
	cgc    conngc.ConnGarbageCollector
}

// close ensures that no new connection can be opened, and all existing ones are closed.
func (d *dialerWrapper) close() {
	d.cgc.CloseAll()
}

func (d *dialerWrapper) Dial() (io.ReadWriteCloser, error) {
	return d.cgc.OpenConn(d.dialer.Dial)
}
