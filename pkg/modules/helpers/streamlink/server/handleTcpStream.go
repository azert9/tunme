package server

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"io"
)

func _handleReceivedDataStream(stream io.ReadWriteCloser, bus *bus) error {

	accepted := bus.sendAcceptedStreamNonBlocking(stream)
	if accepted {
		return nil
	}

	// The stream was not accepted.

	stream.Close()

	return nil
}

func handleTcpStream(stream io.ReadWriteCloser, bus *bus) error {

	var clientHello protocol.ClientHello
	if err := binary.Read(stream, binary.BigEndian, &clientHello); err != nil {
		return err
	}

	switch clientHello.StreamType {
	case protocol.StreamTypeControl:
		return handleControlStream(stream, bus)
	case protocol.StreamTypeConnect:
		return _handleReceivedDataStream(stream, bus)
	case protocol.StreamTypeCallBack:
		bus.sendCallbackStream(stream)
		return nil
	default:
		return fmt.Errorf("invalid stream type in ClientHello")
	}
}
