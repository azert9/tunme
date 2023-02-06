package client

import (
	"encoding/binary"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"log"
	"sync"
	"time"
)

func _runControlStreamOnce(dialer Dialer, bus *bus) error {

	stream, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := binary.Write(stream, binary.BigEndian, protocol.StreamTypeControl); err != nil {
		return err
	}

	for {
		var packetType protocol.ControlPacketType
		if err := binary.Read(stream, binary.BigEndian, &packetType); err != nil {
			return err
		}

		switch packetType {
		case protocol.ControlPacketTypeData:
			// TODO
		case protocol.ControlPacketTypeStreamRequest:
			if !bus.sendAcceptNonBlocking() {
				// the stream was not accepted
				// TODO: concurrency when writing to a control stream
				if err := binary.Write(stream, binary.BigEndian, protocol.ControlPacketTypeStreamRejected); err != nil {
					return err
				}
			}
		}
	}
}

func _runControlStream(dialer Dialer, bus *bus) {

	for {
		err := _runControlStreamOnce(dialer, bus)
		if err != nil {
			// TODO: exit if the tunnel is closed
			log.Print(err)
			time.Sleep(2 * time.Second) // TODO: configure, maybe randomize to avoid all goroutines to dial at once
			continue
		}
	}
}

func runControlStreamsInBackground(wg *sync.WaitGroup, count int, dialer Dialer, bus *bus) {

	for i := 0; i < count; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			_runControlStream(dialer, bus)
		}()
	}
}
