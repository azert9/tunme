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
	bus       bus
}

func NewTunnel(listener net.Listener) modules.Tunnel {

	tun := &tunnel{
		bus: bus{
			acceptedStreamsChan:   make(chan io.ReadWriteCloser),
			receivedPacketsChan:   make(chan []byte),
			outControlPacketsChan: make(chan []byte),
			callBackStreamsChan:   make(chan io.ReadWriteCloser),
		},
	}

	tun.wg.Add(1)
	go func() {
		defer tun.wg.Done()
		acceptLoop(listener, &tun.bus)
	}()

	return tun
}

func (tun *tunnel) Close() (err error) {

	tun.closeOnce.Do(func() {

		// TODO: maybe this is not a good idea, as it can cause writers to panic
		tun.bus.closeAll()

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

	packet, ok := <-tun.bus.receivedPacketsChan
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

	stream, ok := <-tun.bus.acceptedStreamsChan
	if !ok {
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	return stream, nil
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {

	controlPacket := bytes.NewBuffer(make([]byte, 0, 1))
	if err := binary.Write(controlPacket, binary.BigEndian, protocol.ControlPacketTypeStreamRequest); err != nil {
		return nil, err
	}

	// TODO: may panic if the chan is closed
	tun.bus.outControlPacketsChan <- controlPacket.Bytes()

	stream, ok := <-tun.bus.callBackStreamsChan
	if !ok {
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	return stream, nil
}
