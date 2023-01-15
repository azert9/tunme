package tunme_cat

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
	"tunme/pkg/tunme"
)

func forwardStream(in io.Reader, out io.Writer) error {

	buff := make([]byte, 4096)

	for {

		n, readErr := in.Read(buff)

		if _, err := out.Write(buff[:n]); err != nil {
			return err
		}

		if readErr != nil {
			return readErr
		}
	}
}

func closeOrWarn(closer io.Closer) {

	if err := closer.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
}

func cobraMain(_ *cobra.Command, args []string) {

	remote := args[0]

	tunnel, err := tunme.OpenTunnel(remote)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tunnel.OpenStream()
		if err != nil {
			panic(err)
		}
		defer closeOrWarn(conn)

		if err := forwardStream(os.Stdin, conn); err != nil && err != io.EOF {
			panic(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tunnel.AcceptStream()
		if err != nil {
			panic(err)
		}
		defer closeOrWarn(conn)

		if err := forwardStream(conn, os.Stdout); err != nil && err != io.EOF {
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
