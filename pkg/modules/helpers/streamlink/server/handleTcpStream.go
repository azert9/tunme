package server

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"io"
)

func handleTcpStream(stream io.ReadWriteCloser, bus *bus) error {

	var clientHello protocol.ClientHello
	if err := binary.Read(stream, binary.BigEndian, &clientHello); err != nil {
		return err
	}

	switch clientHello.StreamType {
	case protocol.StreamTypeControl:
		return handleControlStream(stream, bus)
	case protocol.StreamTypeConnect:
		bus.acceptedStreamsChan <- stream
		return nil
	case protocol.StreamTypeCallBack:
		bus.callBackStreamsChan <- stream
		return nil
	default:
		return fmt.Errorf("invalid stream type in ClientHello")
	}
}
