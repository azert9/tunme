package helpers

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func makeTestPackets() [][]byte {

	lengths := [4]int{
		0, 1, 158, 1024,
	}

	packets := make([][]byte, len(lengths))

	rnd := rand.New(rand.NewSource(0))

	for i := range packets {

		packets[i] = make([]byte, lengths[i])

		if _, err := rnd.Read(packets[i]); err != nil {
			panic(err)
		}
	}

	return packets[:]
}

var testPackets = makeTestPackets()

func testSinglePacket(t *testing.T, linkType int, roleOrder int, packet []byte) {

	tun1, tun2 := newMockTunPair(linkType, roleOrder)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		received := make([]byte, 3000)
		n, err := tun2.ReceivePacket(received)
		if err != nil {
			t.Logf("error receiving received: %v", err)
			t.Fail()
		}

		if !bytes.Equal(received[:n], packet) {
			t.Logf("received and sent not matching")
			t.Fail()
		}
	}()

	if err := tun1.SendPacket(packet); err != nil {
		t.Logf("error sending packet: %v", err)
		t.Fail()
	}
}

func TestSinglePacket(t *testing.T) {

	for linkType := 0; linkType < 2; linkType++ {

		for roleOrder := 0; roleOrder < 2; roleOrder++ {

			for packetNum := range testPackets {

				t.Run(fmt.Sprintf("link type %d role order %d packet %d", linkType, roleOrder, packetNum), func(t *testing.T) {
					testSinglePacket(t, linkType, roleOrder, testPackets[packetNum])
				})
			}
		}
	}
}
