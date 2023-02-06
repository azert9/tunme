package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"net"
	"sync"
)

type tunnel struct {
	listener  net.Listener
	closeOnce sync.Once
	wg        sync.WaitGroup
	bus       *bus
}

func NewTunnel(listener net.Listener) modules.Tunnel {

	tun := &tunnel{
		listener: listener,
		bus:      newBus(),
	}

	tun.wg.Add(1)
	go func() {
		defer tun.wg.Done()
		acceptLoop(listener, tun.bus)
	}()

	return tun
}

func (tun *tunnel) Close() (err error) {

	tun.closeOnce.Do(func() {

		tun.bus.close()

		err = tun.listener.Close()
		if err != nil {
			return
		}

		tun.wg.Wait()
	})

	return
}

func (tun *tunnel) SendPacket(packet []byte) error {
	// TODO: Send to any control stream, if any available. If none is available, simply return an error.
	//TODO implement me
	panic("implement me")
}

func (tun *tunnel) ReceivePacket(out []byte) (int, error) {

	packet, ok := tun.bus.receiveReceivedPacket()
	if !ok {
		return 0, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	copy(out, packet)

	if len(packet) > len(out) {
		return len(packet), io.ErrShortBuffer
	}

	return len(packet), nil
}

func (tun *tunnel) AcceptStream() (io.ReadWriteCloser, error) {

	stream, ok := tun.bus.receiveAcceptedStream()
	if !ok {
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	// Sending a byte ton confirm that the stream was accepted.

	if _, err := stream.Write([]byte{0}); err != nil {
		stream.Close()
		return nil, err
	}

	//

	return stream, nil
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {

	controlPacket := bytes.NewBuffer(make([]byte, 0, 1))
	if err := binary.Write(controlPacket, binary.BigEndian, protocol.ControlPacketTypeStreamRequest); err != nil {
		return nil, err
	}

	tun.bus.sendOutControlPacket(controlPacket.Bytes())

	stream, ok := tun.bus.receiveCallbackStream()
	if !ok {
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	if stream == nil {
		return nil, modules.ErrStreamRejected
	}

	return stream, nil
}
