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

func cobraMain(_ *cobra.Command, args []string) {

	remote := args[0]
	if flags.Server && flags.Client {
		fmt.Printf("Error: Cannot act both as a server and as a client.")
		os.Exit(1)
	}
	if !flags.Server && !flags.Client {
		// TODO
		panic("hybrid (client and server) mode not implemented")
	}

	tunnel, err := tunme.OpenTunnel(remote)
	if err != nil {
		panic(err)
	}

	var stream io.ReadWriter
	if flags.Server {
		if s, err := tunnel.AcceptStream(); err != nil {
			panic(err)
		} else {
			stream = s
			defer closeOrWarn(s)
		}
	} else {
		if s, err := tunnel.OpenStream(); err != nil {
			panic(err)
		} else {
			stream = s
			defer closeOrWarn(s)
		}
	}

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if _, err := io.Copy(stream, os.Stdin); err != nil && err != io.EOF {
			panic(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if _, err := io.Copy(os.Stdout, stream); err != nil && err != io.EOF {
			panic(err)
		}
	}()
}

var CobraCmd = cobra.Command{
	Use:   "cat REMOTE",
	Short: "Transfer data from standard streams through a tunnel",
	Run:   cobraMain,
	Args:  cobra.ExactArgs(1),
}

var flags struct {
	Server bool
	Client bool
}

func init() {
	CobraCmd.Flags().BoolVar(&flags.Server, "server", false, "If true, will wait for the remote to initiate the stream.")
	CobraCmd.Flags().BoolVar(&flags.Client, "client", false, "If true, will initiate the stream.")
}
