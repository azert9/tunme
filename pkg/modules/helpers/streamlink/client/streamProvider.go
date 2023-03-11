package client

import (
	"context"
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"log"
	"sync"
	"time"
)

type streamProvider struct {
	dialer    Dialer
	available chan io.ReadWriteCloser
}

func newStreamProvider(dialer Dialer) *streamProvider {
	return &streamProvider{
		dialer:    dialer,
		available: make(chan io.ReadWriteCloser), // TODO: configure capacity
	}
}

func (sp *streamProvider) getStream(ctx context.Context) (io.ReadWriteCloser, error) {

	select {
	case conn, ok := <-sp.available:
		if !ok {
			return nil, modules.ErrTunnelClosed
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (sp *streamProvider) loop(wg *sync.WaitGroup, closeChan <-chan struct{}) {

	defer wg.Done()

	cgc := conngc.New()

	for closed := false; !closed; {

		conn, err := cgc.OpenConn(sp.dialer.Dial)
		if err != nil {
			// TODO: detect definitive errors
			// TODO: configure retry delay and timeout
			log.Printf("failed to open connection: %v", err)
			time.Sleep(4 * time.Second)
			continue
		}

		select {
		case sp.available <- conn:
		case <-closeChan:
			closed = true
		}
	}

	close(sp.available)

	cgc.CloseAll()
}
