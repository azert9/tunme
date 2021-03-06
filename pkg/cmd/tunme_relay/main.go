package tunme_relay

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"tunme/pkg/link"
	"tunme/pkg/tunme"
)

// TODO: error handling

func relayDatagrams(tun1 *link.Tunnel, tun2 *link.Tunnel) {

	buff := make([]byte, 100000) // TODO: configure

	for {

		n, _, err := tun1.PacketConn.ReadFrom(buff)
		if err != nil {
			fmt.Printf("error: receiving datagrams: %v", err)
			return
		}

		if _, err := tun2.PacketConn.WriteTo(buff[:n], nil); err != nil {
			fmt.Printf("error: sending datagram: %v", err)
			return
		}
	}
}

func relayStream(conn1 net.Conn, conn2 net.Conn) {

	defer func(conn1 net.Conn) {
		err := conn1.Close()
		if err != nil {
			fmt.Printf("error: closing stream: %v", err)
		}
	}(conn1)

	buff := make([]byte, 4096) // TODO: tune

	for {

		n, readErr := conn1.Read(buff)

		if _, err := conn2.Write(buff[:n]); err != nil {
			fmt.Printf("error: writing to stream: %v", err)
			return
		}

		if readErr == io.EOF {
			return
		} else if readErr != nil {
			fmt.Printf("error: reading from stream: %v", readErr)
			return
		}
	}
}

func relayStreams(tun1 *link.Tunnel, tun2 *link.Tunnel) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	for {

		conn1, err := tun1.StreamListener.Accept()
		if err != nil {
			fmt.Printf("error: accepting connections: %v", err)
			return
		}

		conn2, err := tun2.StreamDialer.Dial()
		if err != nil {
			fmt.Printf("error: opening stream: %v", err)
			return
		}

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			relayStream(conn1, conn2)
		}()

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			relayStream(conn2, conn1)
		}()
	}
}

func relay(tun1 *link.Tunnel, tun2 *link.Tunnel) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		relayStreams(tun1, tun2)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		relayStreams(tun2, tun1)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		relayDatagrams(tun1, tun2)
	}()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		relayDatagrams(tun2, tun1)
	}()
}

func exitBadUsage(program string) {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s REMOTE\n", program)
	os.Exit(1)
}

func Main(program string, args []string) {

	if len(args) != 2 {
		exitBadUsage(program)
	}
	remote1 := args[0]
	remote2 := args[1]

	tun1, err := tunme.OpenTunnel(remote1)
	if err != nil {
		fmt.Printf("error: opening tunnel 1: %v", err)
		os.Exit(1)
	}

	tun2, err := tunme.OpenTunnel(remote2)
	if err != nil {
		fmt.Printf("error: opening tunnel 2: %v", err)
		os.Exit(1)
	}

	relay(tun1, tun2)
}
