package packet_link

import (
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"sync"
	"testing"
	"tunme/pkg/link"
)

var clientTun link.Tunnel
var serverTun link.Tunnel

func init() {

	pipe1 := _mockPacketPipe(make(chan []byte, 100))
	//defer close(pipe1)  // TODO
	pipe2 := _mockPacketPipe(make(chan []byte, 100))
	//defer close(pipe2)  // TODO

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		clientTun = newTunnel(pipe1, pipe2, false)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		serverTun = newTunnel(pipe2, pipe1, true)
	}()
}

func openStream(t *testing.T) (clientStream io.ReadWriteCloser, serverStream io.ReadWriteCloser) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		s, err := clientTun.OpenStream()
		assert.NoError(t, err)
		clientStream = s
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		s, err := serverTun.AcceptStream()
		assert.NoError(t, err)
		serverStream = s
	}()

	return
}

func TestOpenStream(t *testing.T) {

	clientStream, serverStream := openStream(t)
	defer func() {
		assert.NoError(t, clientStream.Close())
		assert.NoError(t, serverStream.Close())
	}()
}

func TestTransferDataClientToServer(t *testing.T) {

	clientStream, serverStream := openStream(t)
	defer func() {
		assert.NoError(t, clientStream.Close())
		assert.NoError(t, serverStream.Close())
	}()

	data := make([]byte, 10000)
	rnd := rand.New(rand.NewSource(0))
	rnd.Read(data)

	// sending data

	n, err := clientStream.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// receiving data

	receiveBuff := make([]byte, len(data))
	_, err = io.ReadFull(serverStream, receiveBuff)
	assert.NoError(t, err)

	assert.Equal(t, data, receiveBuff)
}

func TestTransferDataServerToClient(t *testing.T) {

	clientStream, serverStream := openStream(t)
	defer func() {
		assert.NoError(t, clientStream.Close())
		assert.NoError(t, serverStream.Close())
	}()

	data := make([]byte, 10000)
	rnd := rand.New(rand.NewSource(0))
	rnd.Read(data)

	// sending data

	n, err := serverStream.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// receiving data

	receiveBuff := make([]byte, len(data))
	_, err = io.ReadFull(clientStream, receiveBuff)
	assert.NoError(t, err)

	assert.Equal(t, data, receiveBuff)
}
