package server

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
	"log"
	"math"
)

func (tun *tunnel) controlStreamSendLoop(stream io.ReadWriteCloser) error {

	for done := false; !done; {
		select {
		case packet := <-tun.outControlPacketChan:
			if _, err := stream.Write(packet); err != nil {
				// Failure, putting the packet back in the chan for another stream to pick-up.
				log.Print(err)
				tun.outControlPacketChan <- packet
				done = true
			}
		case <-tun.closeChan:
			done = true
		}
	}

	return nil
}

func (tun *tunnel) controlStreamReceiveDataPacket(stream io.ReadWriteCloser) error {

	defer stream.Close()

	var payloadLen uint32
	if err := binary.Read(stream, binary.BigEndian, &payloadLen); err != nil {
		return err
	}
	if uint64(payloadLen) > math.MaxInt { // TODO: use a configured or fixed value
		return fmt.Errorf("payload too big")
	}
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(stream, payload); err != nil {
		return err
	}

	select {
	case tun.inDataPacketsChan <- payload:
		return nil
	case <-tun.closeChan:
		return modules.ErrTunnelClosed
	}
}

func (tun *tunnel) controlStreamReceiveOne(stream io.ReadWriteCloser) error {

	var packetType protocol.ControlPacketType
	if err := binary.Read(stream, binary.BigEndian, &packetType); err != nil {
		return err
	}

	switch packetType {
	case protocol.ControlPacketTypeData:
		return tun.controlStreamReceiveDataPacket(stream)
	case protocol.ControlPacketTypeStreamRequest:
		return fmt.Errorf("stream requests can only be sent from server to client")
	default:
		return fmt.Errorf("invalid control packet type")
	}
}

func (tun *tunnel) controlStreamReceiveLoop(stream io.ReadWriteCloser) error {

	defer stream.Close()

	for {
		if err := tun.controlStreamReceiveOne(stream); err != nil {
			return err
		}
	}
}

func (tun *tunnel) handleControlStreamInBackground(stream io.ReadWriteCloser) {

	tun.wg.Add(1)
	go func() {
		defer tun.wg.Done()
		if err := tun.controlStreamSendLoop(stream); err != nil {
			log.Printf("error in control stream: %v", err)
		}
	}()

	tun.wg.Add(1)
	go func() {
		defer tun.wg.Done()
		if err := tun.controlStreamReceiveLoop(stream); err != nil {
			log.Printf("error in control stream: %v", err)
		}
	}()
}
