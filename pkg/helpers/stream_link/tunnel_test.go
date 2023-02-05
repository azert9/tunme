package stream_link

import (
	"crypto/rand"
	"fmt"
	"github.com/azert9/tunme/internal/testutil"
	"github.com/azert9/tunme/internal/utils"
	"github.com/azert9/tunme/pkg/link"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func newFakeTunPair() (link.Tunnel, link.Tunnel) {

	var client, server link.Tunnel

	// TODO: redesign testutil

	connectionChan := make(chan testutil.FakeConn)

	client = NewClient(&testutil.FakeDialer{
		OutgoingConnections: connectionChan,
	})

	server = NewServer(&testutil.FakeListener{
		IncomingConnections: connectionChan,
	})

	return client, server
}

func TestTransmitPacket(t *testing.T) {

	client, server := newFakeTunPair()
	data := make([]byte, 10000)
	utils.Must1(io.ReadFull(rand.Reader, data))

	for _, packetSize := range []int{0, 16, 10000} {

		for _, direction := range []string{"client to server", "server to client"} {

			tun1, tun2 := client, server
			if direction == "server to client" {
				tun1, tun2 = server, client
			}

			t.Run(fmt.Sprintf("packet size %d %s", packetSize, direction), func(t *testing.T) {

				err := tun1.SendPacket(data[:packetSize])
				assert.NoError(t, err)

				recvBuff := make([]byte, packetSize*2)
				n, err := tun2.ReceivePacket(recvBuff)
				assert.NoError(t, err)

				assert.Equal(t, data[:packetSize], recvBuff[:n])
			})
		}
	}
}
