package server

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"io"
	"math"
	"sync"
)

func handleControlStream(stream io.ReadWriteCloser, bus *bus) error {

	defer stream.Close()

	var wg sync.WaitGroup
	defer wg.Wait()

	controlStreamClosedChan := make(chan struct{})
	defer close(controlStreamClosedChan)

	// sending control packets

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case controlPacket, ok := <-bus.outControlPacketsChan:
				if !ok {
					return
				}
				if _, err := stream.Write(controlPacket); err != nil {
					// putting the packet back in the channel, as we failed to handle it
					bus.outControlPacketsChan <- controlPacket
					stream.Close()
					return
				}
			case _, _ = <-controlStreamClosedChan:
				return
			}
		}
	}()

	// receiving control packets

	for {

		var packetType protocol.ControlPacketType
		if err := binary.Read(stream, binary.BigEndian, &packetType); err != nil {
			return err
		}

		switch packetType {
		case protocol.ControlPacketTypeData:
			var l uint32
			if err := binary.Read(stream, binary.BigEndian, &l); err != nil {
				return err
			}
			if uint64(l) > math.MaxInt { // TODO: use a configured or fixed value
				return fmt.Errorf("packet too big")
			}
			buff := make([]byte, l)
			if _, err := io.ReadFull(stream, buff); err != nil {
				return err
			}
			bus.receivedPacketsChan <- buff // TODO: this may panic if the chan is closed
		case protocol.ControlPacketTypeStreamRequest:
			return fmt.Errorf("stream requests can only be sent from server to client")
		default:
			return fmt.Errorf("invalid control packet type")
		}
	}
}
