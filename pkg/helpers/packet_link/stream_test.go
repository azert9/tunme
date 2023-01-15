package packet_link

import (
	"fmt"
	"github.com/azert9/tunme/test/assert"
	"io"
	"math/rand"
	"sync"
	"testing"
)

type _mockPacketPipe chan []byte

func (p _mockPacketPipe) SendPacket(packet []byte) error {
	// TODO: drop packets randomly
	p <- packet
	return nil
}

func (p _mockPacketPipe) ReceivePacket(out []byte) (int, error) {

	packet := <-p

	if len(out) < len(packet) {
		return 0, fmt.Errorf("packet too large for buffer")
	}

	copy(out, packet)

	return len(packet), nil
}

func (p _mockPacketPipe) forwardToStream(stream *stream, waitGroup *sync.WaitGroup) {

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			packet, ok := <-p
			if !ok {
				break
			}

			switch packet[0] {
			case byte(packetTypeStreamData):
				dataPacket, err := dataPacketFromBuff(packet)
				if err != nil {
					panic(err)
				}
				stream.handleReceivedDataPacket(dataPacket)
			case byte(packetTypeStreamAck):
				ackPacket, err := ackPacketFromBuff(packet)
				if err != nil {
					panic(err)
				}
				stream.handleReceivedAckPacket(ackPacket)
			}
		}
	}()
}

func TestSendLargeBufferThroughStream(t *testing.T) {

	// Given

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	pipe1 := _mockPacketPipe(make(chan []byte, 100))
	defer close(pipe1)
	pipe2 := _mockPacketPipe(make(chan []byte, 100))
	defer close(pipe2)

	stream1 := newStream(0, pipe1)
	defer stream1.Close()
	stream2 := newStream(0, pipe2)
	defer stream2.Close()

	pipe1.forwardToStream(stream2, &waitGroup)
	pipe2.forwardToStream(stream1, &waitGroup)

	data := make([]byte, 1000)
	rnd := rand.New(rand.NewSource(0))
	rnd.Read(data)

	// When

	sendN, sendErr := stream1.Write(data)

	readBuff := make([]byte, len(data))
	readN, readErr := io.ReadFull(stream2, readBuff)

	// Then

	assert.NoErr(t, sendErr)
	assert.Equal(t, sendN, len(data))

	assert.NoErr(t, readErr)
	assert.Equal(t, readN, len(data))
	assert.SlicesEqual(t, readBuff, data)
}
