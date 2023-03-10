package server

import (
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"log"
	"net"
)

func (tun *tunnel) acceptLoop(listener net.Listener) {

	cgc := conngc.New()

	for {
		conn, err := cgc.OpenConn(func() (io.ReadWriteCloser, error) {
			return listener.Accept()
		})
		if err != nil {
			if !tun.isClosed.Load() {
				log.Printf("error accepting connections: %v", err)
			}
			break
		}

		if err := tun.handleTcpStream(conn); err != nil {
			if err == modules.ErrTunnelClosed {
				break
			}
			log.Print(err)
		}
	}

	// the listener should already be closed when we exit the loop, or be in an error state

	cgc.CloseAll()
}
