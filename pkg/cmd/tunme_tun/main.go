package tunme_tun

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"io"
	"net"
	"os"
	"sync"
	"tunme/pkg/tunme"
)

// TODO: option to automatically add default gateway, or perform NAT

func sendPackets(dev io.Reader, conn net.PacketConn) error {

	buff := make([]byte, 10000) // TODO: configure

	for {

		n, readErr := dev.Read(buff)

		if n > 0 {
			_, err := conn.WriteTo(buff[:n], nil)
			if err != nil {
				return err
			}
		}

		if readErr != nil {
			return readErr
		}
	}
}

func receivePackets(conn net.PacketConn, dev io.Writer) error {

	buff := make([]byte, 10000) // TODO: configure

	for {

		n, _, readErr := conn.ReadFrom(buff)

		if n > 0 {
			_, err := dev.Write(buff[:n])
			if err != nil {
				return err
			}
		}

		if readErr != nil {
			return readErr
		}
	}
}

func exitBadUsage(program string) {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s REMOTE\n", program)
	os.Exit(1)
}

// TODO: help message
type programOptions struct {
	Remote  string `arg:"positional,required"`
	Address string `arg:"-a,--address"`
}

func Main(program string, args []string) {

	var options programOptions
	optParser, err := arg.NewParser(arg.Config{}, &options)
	if err != nil {
		panic(err)
	}
	if err := optParser.Parse(args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n\n", err)
		exitBadUsage(program)
	}

	dev, err := NewTunDevice()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	if options.Address != "" {
		if err := configureInterface(dev.Name(), options.Address); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Created TUN interface \"%s\".\n", dev.Name())

	tunnel, err := tunme.OpenTunnel(options.Remote)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := sendPackets(dev, tunnel.PacketConn); err != nil {
			panic(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := receivePackets(tunnel.PacketConn, dev); err != nil {
			panic(err)
		}
	}()
}
