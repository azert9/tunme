package test

import (
	"github.com/azert9/tunme/pkg/tunme"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"
)

var randomBlockOfData = makeRandomBlockOfData(10000)

func makeRandomBlockOfData(n int) []byte {

	rnd := rand.New(rand.NewSource(0))

	buff := make([]byte, n)
	if _, err := io.ReadFull(rnd, buff); err != nil {
		panic(err)
	}

	return buff
}

func TestFunctional(t *testing.T) {

	if os.Getenv("ENABLE_FUNCTIONAL_TESTS") != "yes" {
		t.SkipNow()
	}

	client, err := tunme.OpenTunnel("tcp-client,localhost:5000")
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		assert.NoError(t, client.Close())
	}()

	server, err := tunme.OpenTunnel("tcp-server,:5000")
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		assert.NoError(t, server.Close())
	}()

	t.Run("client server", func(t *testing.T) {
		tests(t, client, server)
	})

	t.Run("server client", func(t *testing.T) {
		tests(t, server, client)
	})
}

func tests(t *testing.T, tun1 tunme.Tunnel, tun2 tunme.Tunnel) {

	// TODO

	t.Run("opening a stream and transferring data", func(t *testing.T) {

		var wg sync.WaitGroup
		defer wg.Wait()

		// acceptor in background

		wg.Add(1)
		go func() {
			defer wg.Done()

			stream, err := tun2.AcceptStream()
			assert.NoError(t, err)

			buff := make([]byte, len(randomBlockOfData))
			_, err = io.ReadFull(stream, buff)
			assert.NoError(t, err)
			assert.Equal(t, randomBlockOfData, buff)

			n, err := stream.Write(randomBlockOfData)
			assert.NoError(t, err)
			assert.Equal(t, len(randomBlockOfData), n)

			err = stream.Close()
			assert.NoError(t, err)
		}()

		// dialer in foreground

		stream, err := tun1.OpenStream()
		assert.NoError(t, err)

		n, err := stream.Write(randomBlockOfData)
		assert.NoError(t, err)
		assert.Equal(t, len(randomBlockOfData), n)

		buff := make([]byte, len(randomBlockOfData))
		_, err = io.ReadFull(stream, buff)
		assert.NoError(t, err)
		assert.Equal(t, randomBlockOfData, buff)

		n, err = stream.Read(make([]byte, 10))
		assert.Equal(t, io.EOF, err)
		assert.Equal(t, n, 0)

		err = stream.Close()
		assert.NoError(t, err)
	})
}
