package packet_link

import (
	"encoding/binary"
	"fmt"
)

type streamId uint32

type dataPacket struct {
	_buff []byte
}

func dataPacketFromBuff(buff []byte) (dataPacket, error) {

	if len(buff) < 16 {
		return dataPacket{}, fmt.Errorf("truncated dataPacket")
	}

	return dataPacket{
		_buff: buff,
	}, nil
}

func newDataPacket(payloadLen int) dataPacket {
	return dataPacket{
		_buff: make([]byte, 16+payloadLen),
	}
}

func (p dataPacket) getStreamId() streamId {
	return streamId(binary.LittleEndian.Uint32(p._buff[4:]))
}

func (p dataPacket) setStreamId(value streamId) {
	binary.LittleEndian.PutUint32(p._buff[4:], uint32(value))
}

func (p dataPacket) getStreamOffset() uint64 {
	return binary.LittleEndian.Uint64(p._buff[8:])
}

func (p dataPacket) setStreamOffset(value uint64) {
	binary.LittleEndian.PutUint64(p._buff[8:], value)
}

func (p dataPacket) getPayload() []byte {
	return p._buff[16:]
}

type ackPacket struct {
	_buff []byte
}

func newAckPacket() ackPacket {
	return ackPacket{
		_buff: make([]byte, 16),
	}
}

func (p ackPacket) getStreamId() streamId {
	return streamId(binary.LittleEndian.Uint32(p._buff[4:]))
}

func (p ackPacket) setStreamId(value streamId) {
	binary.LittleEndian.PutUint32(p._buff[4:], uint32(value))
}

func (p ackPacket) getStreamOffset() uint64 {
	return binary.LittleEndian.Uint64(p._buff[8:])
}

func (p ackPacket) setStreamOffset(value uint64) {
	binary.LittleEndian.PutUint64(p._buff[8:], value)
}
