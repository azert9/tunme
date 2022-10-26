package tunme_functest

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"tunme/utils"
)

func writeRandom(out io.Writer, quantity int) error {

	buffLenSource := newRandomStream("write-length")
	source := newRandomStream("payload")

	buff := make([]byte, 10000) // TODO: configure

	for quantity > 0 {

		l := utils.Min(quantity, randInt(buffLenSource, len(buff)))

		if _, err := io.ReadFull(source, buff[:l]); err != nil {
			panic(err)
		}

		if _, err := out.Write(buff[:l]); err != nil {
			return err
		}

		quantity -= l
	}

	return nil
}

func readRandom(in io.Reader, quantity int) error {

	buffLenSource := newRandomStream("read-length")
	ref := newRandomStream("payload")

	inBuff := make([]byte, 10000)
	refBuff := make([]byte, 10000)

	for quantity > 0 {

		l := utils.Min(quantity, randInt(buffLenSource, len(inBuff)))

		if _, err := io.ReadFull(ref, refBuff[:l]); err != nil {
			panic(err)
		}

		if _, err := io.ReadFull(in, inBuff[:l]); err != nil {
			return err
		}

		if !bytes.Equal(inBuff, refBuff) {
			return fmt.Errorf("received stream does not match reference")
		}

		quantity -= l
	}

	return nil
}

func testStream(conn io.ReadWriteCloser, quantity int) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	// Writing

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if err := writeRandom(conn, quantity); err != nil && err != io.EOF {
			panic(err)
		}
	}()

	// Reading

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if err := readRandom(conn, quantity); err != nil && err != io.EOF {
			panic(err)
		}
	}()
}
