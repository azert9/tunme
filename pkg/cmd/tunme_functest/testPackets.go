package tunme_functest

import (
	"bytes"
	"fmt"
	"github.com/azert9/tunme/internal/utils"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"sync"
)

func sendPacketsRandom(tun modules.PacketSender, quantity int) error {

	source := newRandomStream("payload")

	buff := make([]byte, 10000) // TODO: configure

	for quantity > 0 {

		l := utils.Min(quantity, randInt(source, len(buff)))

		if _, err := io.ReadFull(source, buff[:l]); err != nil {
			panic(err)
		}

		if err := tun.SendPacket(buff[:l]); err != nil {
			return err
		}

		quantity -= l
	}

	return nil
}

func receivePacketsRandom(in modules.PacketReceiver, quantity int) error {

	source := newRandomStream("payload")

	inBuff := make([]byte, 10000)
	refBuff := make([]byte, 10000)

	for quantity > 0 {

		l := utils.Min(quantity, randInt(source, len(inBuff)))

		if _, err := io.ReadFull(source, refBuff[:l]); err != nil {
			panic(err)
		}

		if _, err := in.ReceivePacket(inBuff[:l]); err != nil {
			return err
		}

		if !bytes.Equal(inBuff, refBuff) {
			return fmt.Errorf("received datagram does not match reference")
		}

		quantity -= l
	}

	return nil
}

func testDatagrams(tun modules.Tunnel, quantity int) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	// Writing

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if err := sendPacketsRandom(tun, quantity); err != nil && err != io.EOF {
			panic(err)
		}
	}()

	// Reading

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if err := receivePacketsRandom(tun, quantity); err != nil && err != io.EOF {
			panic(err)
		}
	}()
}
