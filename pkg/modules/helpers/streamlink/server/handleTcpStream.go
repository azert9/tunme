package server

import (
	"encoding/binary"
	"fmt"
	"github.com/azert9/tunme/internal/streamlink/protocol"
	"github.com/azert9/tunme/pkg/modules"
	"io"
)

func (tun *tunnel) handleTcpStream(stream io.ReadWriteCloser) error {

	var clientHello protocol.ClientHello
	if err := binary.Read(stream, binary.BigEndian, &clientHello); err != nil {
		return err
	}

	switch clientHello.StreamType {
	case protocol.StreamTypeControl:
		tun.handleControlStreamInBackground(stream)
	case protocol.StreamTypeConnect:
		select {
		case tun.inDataStreamChan <- stream:
		case <-tun.closeChan:
			return modules.ErrTunnelClosed
		}
	case protocol.StreamTypeCallBack:
		select {
		case tun.inCallBackStreamChan <- stream:
		case <-tun.closeChan:
			return modules.ErrTunnelClosed
		}
	default:
		return fmt.Errorf("invalid stream type in ClientHello")
	}

	return nil
}
