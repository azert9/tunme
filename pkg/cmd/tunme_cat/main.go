package tunme_cat

import (
	"fmt"
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

func exitBadUsage(program string) {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s REMOTE REMOTE\n", program)
	os.Exit(1)
}

func Main(program string, args []string) {

	if len(args) != 1 {
		exitBadUsage(program)
	}
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

		conn, err := tunnel.StreamDialer.Dial()
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

		conn, err := tunnel.StreamListener.Accept()
		if err != nil {
			panic(err)
		}
		defer closeOrWarn(conn)

		if err := forwardStream(conn, os.Stdout); err != nil && err != io.EOF {
			panic(err)
		}
	}()
}
