package tunme_tun

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/tunme"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
)

// TODO: option to automatically add default gateway, or perform NAT

func sendPackets(dev io.Reader, tun modules.PacketSender) error {

	buff := make([]byte, 10000) // TODO: configure

	for {

		n, readErr := dev.Read(buff)

		if n > 0 {
			err := tun.SendPacket(buff[:n])
			if err != nil {
				return err
			}
		}

		if readErr != nil {
			return readErr
		}
	}
}

func receivePackets(tun modules.PacketReceiver, dev io.Writer) error {

	buff := make([]byte, 10000) // TODO: configure

	for {

		n, readErr := tun.ReceivePacket(buff)

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

func cobraMain(_ *cobra.Command, args []string) {

	remote := args[0]

	dev, err := NewTunDevice()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	for _, addr := range flags.Addresses {
		if err := addIpAddressToInterface(dev.Name(), addr); err != nil {
			panic(err)
		}
		fmt.Printf("addr: %s\n", addr)
	}

	fmt.Printf("Created TUN interface \"%s\".\n", dev.Name())

	tunnel, err := tunme.OpenTunnel(remote)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := sendPackets(dev, tunnel); err != nil {
			panic(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := receivePackets(tunnel, dev); err != nil {
			panic(err)
		}
	}()
}

var CobraCmd = cobra.Command{
	Use:   "tun REMOTE",
	Short: "Create a virtual network interface",
	Run:   cobraMain,
	Args:  cobra.ExactArgs(1),
}

var flags struct {
	Addresses []string
}

func init() {
	CobraCmd.Flags().StringArrayVarP(&flags.Addresses, "address", "a", nil, "Assign an address to the interface. Can be repeated.")
}
