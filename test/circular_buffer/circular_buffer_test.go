package circular_buffer

import (
	"testing"
	"tunme/internal/circular_buffer"
	"tunme/test/assert"
)

var sampleData1 = []byte{0x34, 0x87, 0x01, 0xf9, 0x56, 0x32, 0x11}

func TestWriteLessThanCapacity(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) * 2)

	// When

	writeN, writeErr := buff.Write(sampleData1)

	// Then

	assert.NoErr(t, writeErr)

	assert.Equal(t, writeN, len(sampleData1))
}

func TestReadAvailable(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(sampleData1)

	// When

	readBuff := make([]byte, len(sampleData1)*2)
	n, err := buff.Read(readBuff)

	// Then

	assert.NoErr(t, err)

	assert.SlicesEqual(t, readBuff[:n], sampleData1)
}

func TestReadAvailableInTwoChunks(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(sampleData1)
	_, _ = buff.Read(make([]byte, len(sampleData1)/2))

	// When

	readBuff := make([]byte, len(sampleData1)*2)
	readN, readErr := buff.Read(readBuff)

	// Then

	assert.NoErr(t, readErr)

	assert.SlicesEqual(t, readBuff[:readN], sampleData1[len(sampleData1)/2:])
}

func TestWriteInTwoChunks(t *testing.T) {

	// Given

	chunk1 := sampleData1[:len(sampleData1)/2]
	chunk2 := sampleData1[len(sampleData1)/2:]

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(chunk1)

	// When

	writeN, writeErr := buff.Write(chunk2)

	// Then

	assert.NoErr(t, writeErr)

	assert.Equal(t, writeN, len(chunk2))

	readBuff := make([]byte, len(sampleData1)*2)
	n, err := buff.Read(readBuff)
	assert.NoErr(t, err)

	assert.SlicesEqual(t, readBuff[:n], sampleData1)
}

func TestWriteMoreThanCapacity(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(10)

	// When

	n, err := buff.Write(make([]byte, buff.Capacity()+1))

	// Then

	if err != circular_buffer.ErrNotEnoughSpaceInBuffer {
		t.Logf("unexpected error (or missing error): %v", err)
		t.Fail()
	}

	assert.Equal(t, n, 0)
}

func TestWriteAroundTheBuffer(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) + 4)
	_, _ = buff.Write(make([]byte, len(sampleData1)))
	_, _ = buff.Read(make([]byte, len(sampleData1)))

	// When

	writeN, writeErr := buff.Write(sampleData1)

	// Then

	assert.NoErr(t, writeErr)
	assert.Equal(t, writeN, len(sampleData1))
}

func TestReadAroundTheBuffer(t *testing.T) {

	// Given

	buff := circular_buffer.NewCircularBuffer(len(sampleData1) + 4)
	_, _ = buff.Write(make([]byte, len(sampleData1)))
	_, _ = buff.Read(make([]byte, len(sampleData1)))
	_, _ = buff.Write(sampleData1)

	// When

	readBuff := make([]byte, len(sampleData1)*2)
	readN, readErr := buff.Read(readBuff)

	// Then

	assert.NoErr(t, readErr)
	assert.SlicesEqual(t, readBuff[:readN], sampleData1)
}
