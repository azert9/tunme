package client

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/conngc"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"sync"
)

type Dialer interface {
	Dial() (io.ReadWriteCloser, error)
}

type tunnel struct {
	wg     sync.WaitGroup
	dialer *dialerWrapper
	bus    *bus
}

func NewTunnel(dialer Dialer) modules.Tunnel {

	tun := &tunnel{
		dialer: &dialerWrapper{
			dialer: dialer,
			cgc:    conngc.New(),
		},
		bus: newBus(),
	}

	runControlStreamsInBackground(&tun.wg, 4, dialer, tun.bus) // TODO: configure count

	return tun
}

func (tun *tunnel) Close() error {
	tun.bus.close()
	tun.dialer.close()
	return nil
}

func (tun *tunnel) SendPacket(packet []byte) error {
	//TODO implement me
	panic("implement me")
}

func (tun *tunnel) ReceivePacket(out []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (tun *tunnel) AcceptStream() (io.ReadWriteCloser, error) {

	ok := tun.bus.receiveAccept()
	if !ok {
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	stream, err := tun.dialer.Dial()
	if err != nil {
		return nil, err
	}

	if err := binary.Write(stream, binary.BigEndian, protocol.StreamTypeCallBack); err != nil {
		stream.Close()
		return nil, err
	}

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

	// The server should send a single byte to confirm that the stream was accepted.

	var buff [1]byte
	_, err = conn.Read(buff[:])
	if err != nil {
		if err == io.EOF {
			return nil, modules.ErrStreamRejected
		}
		return nil, err
	}

	//

	return conn, nil
}
