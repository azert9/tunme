package packet_link

import (
	"io"
	"testing"
	"tunme/internal/circular_buffer"
	"tunme/test/assert"
)

var _testPayload = []byte{0x54, 0x98, 0xa6, 0x00, 0x2f, 0xff, 0xb8}

func _makeTestPacketWithPayload(offset uint64, payload []byte) dataPacket {

	packet := newDataPacket(len(payload))
	packet.setStreamOffset(offset)
	copy(packet.getPayload(), payload)

	return packet
}

func _makeTestPacket(offset uint64) dataPacket {
	return _makeTestPacketWithPayload(offset, _testPayload)
}

type _mockAckInfo struct {
	streamId streamId
	offset   uint64
}

type _mockAckSender struct {
	sent []_mockAckInfo
}

func (s *_mockAckSender) sendAck(streamId streamId, offset uint64) error {

	s.sent = append(s.sent, _mockAckInfo{
		streamId: streamId,
		offset:   offset,
	})

	return nil
}

func _newMockAckSender() *_mockAckSender {
	return &_mockAckSender{}
}

func TestReceivePacketNoError(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	packet := _makeTestPacket(0)

	// When

	receiveErr := stream.handleReceivedDataPacket(packet)

	// Then

	assert.NoErr(t, receiveErr)
}

func TestReadDataJustReceived(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacket(0))

	// When

	readBuff := make([]byte, len(_testPayload))
	readN, readErr := io.ReadFull(stream.buff, readBuff)

	// Then

	assert.NoErr(t, readErr)
	assert.SlicesEqual(t, _testPayload, readBuff[:readN])
}

func TestStreamDropsUnorderedPacket(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	packet := _makeTestPacket(1)

	// When

	receiveErr := stream.handleReceivedDataPacket(packet)

	// Then

	if receiveErr == nil {
		t.Logf("error is nil")
		t.Fail()
	}
}

func TestReceiveTwoChunks(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacketWithPayload(0, _testPayload[:3]))
	packet2 := _makeTestPacketWithPayload(3, _testPayload[3:])

	// When

	receiveErr := stream.handleReceivedDataPacket(packet2)

	// Then

	assert.NoErr(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoErr(t, err)
	assert.SlicesEqual(t, _testPayload, readBuff[:n])
}

func TestReceiveOverlappingChunks(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacketWithPayload(0, _testPayload[:3]))
	packet2 := _makeTestPacketWithPayload(2, _testPayload[2:])

	// When

	receiveErr := stream.handleReceivedDataPacket(packet2)

	// Then

	assert.NoErr(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoErr(t, err)
	assert.SlicesEqual(t, _testPayload, readBuff[:n])
}

func TestReceiveChunkInThePast(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacket(0))
	packet2 := _makeTestPacketWithPayload(2, _testPayload[2:3])

	// When

	receiveErr := stream.handleReceivedDataPacket(packet2)

	// Then

	assert.NoErr(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoErr(t, err)
	assert.SlicesEqual(t, _testPayload, readBuff[:n])
}

func TestAckAfterPacketReceived(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	packet := _makeTestPacket(0)

	// When

	_ = stream.handleReceivedDataPacket(packet)

	// Then

	assert.Equal(t, len(mockAckSender.sent), 1)
	assert.Equal(t, mockAckSender.sent[0].streamId, 0) // TODO
	assert.Equal(t, mockAckSender.sent[0].offset, uint64(len(_testPayload)))
}

func TestAckAfterRetransmission(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacket(0))

	// When

	_ = stream.handleReceivedDataPacket(_makeTestPacketWithPayload(1, _testPayload[1:3]))

	// Then

	assert.Equal(t, len(mockAckSender.sent), 2)
	assert.Equal(t, mockAckSender.sent[1].streamId, 0) // TODO
	assert.Equal(t, mockAckSender.sent[1].offset, uint64(len(_testPayload)))
}

func TestAckNotSentAfterUnorderedPacketInTheFuture(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))

	// When

	_ = stream.handleReceivedDataPacket(_makeTestPacket(1))

	// Then

	assert.Equal(t, len(mockAckSender.sent), 0)
}

// TODO: test receive window full
