//go:build !without_cat_app

package tunme_cat

import (
	"fmt"
	"github.com/azert9/tunme/pkg/tunme"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
)

func closeOrWarn(closer io.Closer) {

	if err := closer.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
}

func send(tun tunme.Tunnel) error {

	stream, err := tun.OpenStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = io.Copy(stream, os.Stdin)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func receive(tun tunme.Tunnel) error {

	stream, err := tun.AcceptStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = io.Copy(os.Stdout, stream)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func cobraMain(_ *cobra.Command, args []string) {

	dosSend := true
	dorReceive := true
	if flags.Send || flags.Receive {
		dosSend = flags.Send
		dorReceive = flags.Receive
	}

	remote := args[0]

	tun, err := tunme.OpenTunnel(remote)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	if dosSend {

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			if err := send(tun); err != nil {
				panic(err)
			}
		}()
	}

	if dorReceive {

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			if err := receive(tun); err != nil {
				panic(err)
			}
		}()
	}
}

var flags struct {
	Send    bool
	Receive bool
}

func RegisterCmd(parentCmd *cobra.Command) {

	cmd := cobra.Command{
		Use:   "cat REMOTE",
		Short: "Transfer data from standard streams through a tunnel",
		Long:  "By default, the communication is bidirectional. Use --send and --receive to control this behavior.",
		Run:   cobraMain,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().BoolVar(&flags.Send, "send", false, "Send data to the peer.")
	cmd.Flags().BoolVar(&flags.Receive, "receive", false, "Receive data from the peer")

	parentCmd.AddCommand(&cmd)
}
