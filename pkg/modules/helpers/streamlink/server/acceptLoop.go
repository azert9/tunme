package server

import (
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"io"
	"log"
	"net"
	"sync"
)

func acceptLoop(listener net.Listener, bus *bus) {

	cgc := conngc.New()

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		conn, err := cgc.OpenConn(func() (io.ReadWriteCloser, error) {
			return listener.Accept()
		})
		if err != nil {
			log.Println(err)
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := handleTcpStream(conn, bus); err != nil {
				log.Print(err)
			}
		}()
	}

	// the listener should already be closed when we exit the loop, or be in an error state

	cgc.CloseAll()
}
