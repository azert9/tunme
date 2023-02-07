package packet_link

import (
	"github.com/azert9/tunme/internal/circular_buffer"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
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

	assert.NoError(t, receiveErr)
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

	assert.NoError(t, readErr)
	assert.Equal(t, readBuff[:readN], _testPayload)
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

	assert.NoError(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoError(t, err)
	assert.Equal(t, readBuff[:n], _testPayload)
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

	assert.NoError(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoError(t, err)
	assert.Equal(t, readBuff[:n], _testPayload)
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

	assert.NoError(t, receiveErr)

	readBuff := make([]byte, len(_testPayload))
	n, err := stream.buff.Read(readBuff)
	assert.NoError(t, err)
	assert.Equal(t, readBuff[:n], _testPayload)
}

func TestAckAfterPacketReceived(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	packet := _makeTestPacket(0)

	// When

	_ = stream.handleReceivedDataPacket(packet)

	// Then

	assert.Equal(t, 1, len(mockAckSender.sent))
	assert.Equal(t, streamId(0), mockAckSender.sent[0].streamId) // TODO
	assert.Equal(t, uint64(len(_testPayload)), mockAckSender.sent[0].offset)
}

func TestAckAfterRetransmission(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))
	_ = stream.handleReceivedDataPacket(_makeTestPacket(0))

	// When

	_ = stream.handleReceivedDataPacket(_makeTestPacketWithPayload(1, _testPayload[1:3]))

	// Then

	assert.Equal(t, 2, len(mockAckSender.sent))
	assert.Equal(t, streamId(0), mockAckSender.sent[1].streamId) // TODO
	assert.Equal(t, uint64(len(_testPayload)), mockAckSender.sent[1].offset)
}

func TestAckNotSentAfterUnorderedPacketInTheFuture(t *testing.T) {

	// Given

	mockAckSender := _newMockAckSender()
	stream := newStreamReceiver(0, mockAckSender, circular_buffer.NewCircularBuffer(200))

	// When

	_ = stream.handleReceivedDataPacket(_makeTestPacket(1))

	// Then

	assert.Equal(t, 0, len(mockAckSender.sent))
}

// TODO: test receive window full
