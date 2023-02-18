package server

import (
	"bytes"
	"encoding/binary"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type tunnel struct {
	listener             net.Listener
	wg                   sync.WaitGroup
	isClosed             atomic.Bool
	closeChan            chan struct{}
	outControlPacketChan chan []byte
	inDataPacketsChan    chan []byte
	inDataStreamChan     chan io.ReadWriteCloser
	inCallBackStreamChan chan io.ReadWriteCloser
}

func NewTunnel(listener net.Listener) modules.Tunnel {

	// TODO: configure
	packetBacklog := 16
	streamBacklog := 16

	tun := &tunnel{
		listener:             listener,
		closeChan:            make(chan struct{}),
		outControlPacketChan: make(chan []byte, 32),
		inDataPacketsChan:    make(chan []byte, packetBacklog),
		inDataStreamChan:     make(chan io.ReadWriteCloser, streamBacklog),
		inCallBackStreamChan: make(chan io.ReadWriteCloser, 32),
	}

	tun.wg.Add(1)
	go func() {
		defer tun.wg.Done()
		tun.acceptLoop(listener)
	}()

	return tun
}

func (tun *tunnel) Close() (err error) {

	if tun.isClosed.Swap(true) {
		return nil
	}

	close(tun.closeChan)

	err = tun.listener.Close()
	if err != nil {
		return
	}

	tun.wg.Wait()

	return nil
}

func (tun *tunnel) SendPacket(packet []byte) error {

	// TODO: should ensure that the packet channel cannot be used once the tunnel is closed
	select {
	case tun.outControlPacketChan <- protocol.BuildDataControlPacket(packet):
		return nil
	case <-tun.closeChan:
		return modules.ErrTunnelClosed
	}
}

func (tun *tunnel) ReceivePacket(out []byte) (int, error) {

	var packet []byte
	select {
	case packet = <-tun.inDataPacketsChan:
	case <-tun.closeChan:
		return 0, modules.ErrTunnelClosed
	}

	copy(out, packet)

	if len(packet) > len(out) {
		return len(packet), io.ErrShortBuffer
	}

	return len(packet), nil
}

func (tun *tunnel) AcceptStream() (io.ReadWriteCloser, error) {

	select {
	case stream := <-tun.inDataStreamChan:
		return stream, nil
	case <-tun.closeChan:
		return nil, modules.ErrTunnelClosed
	}
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {

	controlPacket := bytes.NewBuffer(make([]byte, 0, 1))
	if err := binary.Write(controlPacket, binary.BigEndian, protocol.ControlPacketTypeStreamRequest); err != nil {
		return nil, err
	}

	select {
	case tun.outControlPacketChan <- controlPacket.Bytes():
	case <-tun.closeChan:
		return nil, modules.ErrTunnelClosed
	}

	select {
	case stream := <-tun.inCallBackStreamChan:
		return stream, nil
	case <-tun.closeChan:
		return nil, modules.ErrTunnelClosed
	}
}
