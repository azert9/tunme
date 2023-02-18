package protocol

import (
	"bytes"
	"encoding/binary"
)

type StreamType byte

const (
	StreamTypeControl StreamType = 0
	StreamTypeConnect StreamType = 1
	// StreamTypeCallBack is for when the client opens a connection at the request of the server.
	StreamTypeCallBack StreamType = 2
)

// ClientHello is sent by the client upon opening a connection.
type ClientHello struct {
	StreamType StreamType
}

type ControlPacketType byte

const (
	ControlPacketTypeData          ControlPacketType = 0
	ControlPacketTypeStreamRequest ControlPacketType = 1
)

func BuildDataControlPacket(data []byte) []byte {

	packet := bytes.NewBuffer(make([]byte, 0, 5+len(data)))

	if err := binary.Write(packet, binary.BigEndian, ControlPacketTypeData); err != nil {
		panic(err)
	}

	if err := binary.Write(packet, binary.BigEndian, uint32(len(data))); err != nil {
		panic(err)
	}

	if _, err := packet.Write(data); err != nil {
		panic(err)
	}

	return packet.Bytes()
}
