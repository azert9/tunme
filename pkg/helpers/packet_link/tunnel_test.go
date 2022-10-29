package packet_link

import (
	"math/rand"
	"sync"
	"testing"
	"tunme/test/assert"
)

func TestOpenStream(t *testing.T) {

	// Given

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	pipe1 := _mockPacketPipe(make(chan []byte, 100))
	defer close(pipe1)
	pipe2 := _mockPacketPipe(make(chan []byte, 100))
	defer close(pipe2)

	clientTun := newTunnel(pipe1, pipe2, false)
	serverTun := newTunnel(pipe2, pipe1, true)

	data := make([]byte, 1000)
	rnd := rand.New(rand.NewSource(0))
	rnd.Read(data)

	// When

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		// TODO
		_, _ = clientTun.OpenStream()
	}()

	_, acceptErr := serverTun.AcceptStream()

	// Then

	assert.NoErr(t, acceptErr)
}
