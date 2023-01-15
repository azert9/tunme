package tunme_functest

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/pkg/tunme"
	"io"
	"os"
	"sync"
)

func randInt(reader io.Reader, upperBound int) int {

	upperBoundBeforeModulo := uint32((1 << 32) / upperBound * upperBound)

	for {

		var val uint32
		if err := binary.Read(reader, binary.BigEndian, &val); err != nil {
			panic(err)
		}

		if val >= upperBoundBeforeModulo {
			// avoiding modulo bias
			continue
		}

		return int(val) % upperBound
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

	quantity := 1000000

	for i := 0; i < 3; i++ {

		// Writing Stream

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			conn, err := tunnel.OpenStream()
			if err != nil {
				panic(err)
			}
			defer closeOrWarn(conn)

			testStream(conn, quantity)
		}()

		// Reading Stream

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			conn, err := tunnel.AcceptStream()
			if err != nil {
				panic(err)
			}
			defer closeOrWarn(conn)

			testStream(conn, quantity)
		}()
	}

	// Datagrams

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		testDatagrams(tunnel, quantity)
	}()
}
