package tunme_tcp

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/tunme"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
	"os"
)

func forwardStreamInBackground(a io.ReadWriteCloser, b io.ReadWriteCloser) {

	// TODO: ensure no goroutine is left behind

	go func() {
		defer a.Close()
		if _, err := io.Copy(a, b); err != nil && err != io.EOF {
			log.Print(err)
		}
	}()

	go func() {
		defer b.Close()
		if _, err := io.Copy(b, a); err != nil && err != io.EOF {
			log.Print(err)
		}
	}()
}

func doServer(tun modules.Tunnel, address string) error {

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		tunStream, err := tun.OpenStream()
		if err != nil {
			conn.Close()
			return err
		}

		forwardStreamInBackground(tunStream, conn)
	}
}

func doClient(tun modules.Tunnel, address string) error {

	for {

		tunStream, err := tun.AcceptStream()
		if err != nil {
			return err
		}

		clientConn, err := net.Dial("tcp", address)
		if err != nil {
			tunStream.Close()
			// TODO: temporary error
			return err
		}

		forwardStreamInBackground(tunStream, clientConn)
	}
}

func cobraMain(_ *cobra.Command, args []string) {

	var isServer bool
	if args[0] == "listen" {
		isServer = true
	} else if args[0] == "connect" {
		isServer = false
	} else {
		fmt.Printf("invalid argument: must specify either \"listen\" or \"connect\", not %q\n", args[0])
		os.Exit(1)
	}
	address := args[1]
	remote := args[2]

	tun, err := tunme.OpenTunnel(remote)
	if err != nil {
		fmt.Printf("failed to create tunnel: %v", err)
		os.Exit(1)
	}
	defer tun.Close()

	// TODO: open the tunnel in a loop

	if isServer {
		if err := doServer(tun, address); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if err := doClient(tun, address); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

var CobraCmd = cobra.Command{
	Use:   "tcp listen|connect ADDRESS REMOTE",
	Short: "Forward a TCP port",
	Run:   cobraMain,
	Args:  cobra.ExactArgs(3),
}
