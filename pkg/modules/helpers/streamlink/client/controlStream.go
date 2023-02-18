package client

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"log"
	"sync"
	"time"
)

func (tun *tunnel) controlStreamReceiveDataPacket(stream io.ReadWriteCloser) error {

	var payloadLen uint32
	if err := binary.Read(stream, binary.BigEndian, &payloadLen); err != nil {
		return err
	}
	if payloadLen > 5000 { // TODO
		return fmt.Errorf("packet too big")
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(stream, payload); err != nil {
		return err
	}

	select {
	case tun.inDataPacketChan <- payload:
	case <-tun.closeChan:
		return modules.ErrTunnelClosed
	}

	return nil
}

func (tun *tunnel) controlStreamReceiveOne(stream io.ReadWriteCloser) error {

	var packetType protocol.ControlPacketType
	if err := binary.Read(stream, binary.BigEndian, &packetType); err != nil {
		return err
	}

	switch packetType {
	case protocol.ControlPacketTypeData:
		if err := tun.controlStreamReceiveDataPacket(stream); err != nil {
			return err
		}
	case protocol.ControlPacketTypeStreamRequest:
		select {
		case tun.inStreamRequestChan <- streamRequest{}:
		case <-tun.closeChan:
			return modules.ErrTunnelClosed
		}
	default:
		return fmt.Errorf("invalid control packet type")
	}

	return nil
}

func (tun *tunnel) controlStreamReceiveLoop(stream io.ReadWriteCloser) error {

	defer stream.Close()

	for {
		if err := tun.controlStreamReceiveOne(stream); err != nil {
			return err
		}
	}
}

func (tun *tunnel) controlStreamSendLoop(stream io.ReadWriteCloser) error {

	defer stream.Close()

	if err := binary.Write(stream, binary.BigEndian, protocol.StreamTypeControl); err != nil {
		return err
	}

	for {
		select {
		case packet := <-tun.outControlPacketChan:
			if _, err := stream.Write(packet); err != nil {
				// Failure, putting the packet back in the chan for another stream to pick-up.
				tun.outControlPacketChan <- packet
				return err
			}
		case <-tun.closeChan:
			return modules.ErrTunnelClosed
		}
	}
}

func (tun *tunnel) runControlStreamOnce() error {

	var wg sync.WaitGroup
	defer wg.Wait()

	stream, err := tun.dialer.Dial()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tun.controlStreamSendLoop(stream); err != nil {
			if !tun.isClosed.Load() {
				log.Printf("error in control stream: %v", err)
			}
		}
	}()

	if err := tun.controlStreamReceiveLoop(stream); err != nil {
		if !tun.isClosed.Load() {
			log.Printf("error in control stream: %v", err)
		}
	}

	return nil
}

func (tun *tunnel) runControlStream() {

	for !tun.isClosed.Load() {

		if err := tun.runControlStreamOnce(); err != nil {
			if err == modules.ErrTunnelClosed {
				return
			}
			log.Print(err)
		}

		if tun.isClosed.Load() {
			break
		}

		time.Sleep(2 * time.Second) // TODO: configure
	}
}
