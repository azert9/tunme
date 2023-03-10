package client

import (
	"encoding/binary"
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"sync"
	"sync/atomic"
)

type Dialer interface {
	Dial() (io.ReadWriteCloser, error)
}

type streamRequest struct {
}

type tunnel struct {
	wg                   sync.WaitGroup
	dialer               *dialerWrapper
	isClosed             atomic.Bool
	closeChan            chan struct{}
	outControlPacketChan chan []byte
	inDataPacketChan     chan []byte
	inStreamRequestChan  chan streamRequest
}

func NewTunnel(dialer Dialer) modules.Tunnel {

	tun := &tunnel{
		dialer: &dialerWrapper{
			dialer: dialer,
			cgc:    conngc.New(),
		},
		// TODO: channel sizes
		closeChan:            make(chan struct{}),
		outControlPacketChan: make(chan []byte),
		inDataPacketChan:     make(chan []byte),
		inStreamRequestChan:  make(chan streamRequest),
	}

	for i := 0; i < 4; i++ {
		tun.wg.Add(1)
		go func() {
			defer tun.wg.Done()
			tun.runControlStream()
		}()
	}

	return tun
}

func (tun *tunnel) Close() error {

	if tun.isClosed.Swap(true) {
		return nil
	}

	tun.dialer.close()

	close(tun.closeChan)

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
	case packet = <-tun.inDataPacketChan:
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
	case <-tun.inStreamRequestChan:
	case <-tun.closeChan:
		return nil, modules.ErrTunnelClosed
	}

	streamSetupOk := false

	stream, err := tun.dialer.Dial()
	if err != nil {
		// TODO: Retry? Send a message?
		return nil, err
	}
	defer func() {
		if !streamSetupOk {
			stream.Close()
		}
	}()

	if err := binary.Write(stream, binary.BigEndian, protocol.StreamTypeCallBack); err != nil {
		return nil, err
	}

	streamSetupOk = true
	return stream, nil
}

func (tun *tunnel) OpenStream() (io.ReadWriteCloser, error) {

	conn, err := tun.dialer.Dial()
	if err != nil {
		return nil, err
	}

	streamType := protocol.StreamTypeConnect
	if err := binary.Write(conn, binary.BigEndian, streamType); err != nil {
		return nil, err
	}

	return conn, nil
}
