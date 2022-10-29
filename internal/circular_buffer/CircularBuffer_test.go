package circular_buffer

import (
	"io"
	"sync"
	"testing"
	"time"
	"tunme/test/assert"
)

var sampleData1 = []byte{0x34, 0x87, 0x01, 0xf9, 0x56, 0x32, 0x11}

func TestWriteLessThanCapacity(t *testing.T) {

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)

	// When

	writeN, writeErr := buff.Write(sampleData1)

	// Then

	assert.NoErr(t, writeErr)
	assert.Equal(t, buff.Len(), len(sampleData1))
	assert.Equal(t, writeN, len(sampleData1))
}

func TestReadAvailable(t *testing.T) {

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(sampleData1)

	// When

	readBuff := make([]byte, len(sampleData1)*2)
	n, err := buff.Read(readBuff)

	// Then

	assert.NoErr(t, err)

	assert.SlicesEqual(t, readBuff[:n], sampleData1)
}

func TestPeekAvailable(t *testing.T) {

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(sampleData1)

	// When

	readBuff := make([]byte, len(sampleData1)*2)
	n, err := buff.Peek(readBuff)

	// Then

	assert.NoErr(t, err)

	assert.SlicesEqual(t, readBuff[:n], sampleData1)
}

func TestReadAvailableInTwoChunks(t *testing.T) {

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)
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

	buff := NewCircularBuffer(len(sampleData1) * 2)
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

	buff := NewCircularBuffer(10)

	// When

	n, err := buff.Write(make([]byte, buff.Capacity()+1))

	// Then

	if err != ErrNotEnoughSpaceInBuffer {
		t.Logf("unexpected error (or missing error): %v", err)
		t.Fail()
	}

	assert.Equal(t, n, 0)
}

func TestWriteAroundTheBuffer(t *testing.T) {

	// Given

	buff := NewCircularBuffer(len(sampleData1) + 4)
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

	buff := NewCircularBuffer(len(sampleData1) + 4)
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

func TestBlockingRead(t *testing.T) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)

	// When

	readBuff := make([]byte, len(sampleData1)*2)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		time.Sleep(200 * time.Millisecond)
		_, _ = buff.Write(sampleData1)
	}()

	readN, readErr := buff.Read(readBuff)

	// Then

	assert.NoErr(t, readErr)

	assert.NotEqual(t, readN, 0)

	assert.SlicesEqual(t, readBuff[:readN], sampleData1[:readN])
}

func TestPeekDoesNotConsumeData(t *testing.T) {

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	// Given

	buff := NewCircularBuffer(len(sampleData1) * 2)
	_, _ = buff.Write(sampleData1)

	// When

	_, _ = buff.Peek(make([]byte, len(sampleData1)))

	// Then

	n, err := io.ReadFull(buff, make([]byte, len(sampleData1)))
	assert.NoErr(t, err)
	assert.Equal(t, n, len(sampleData1))
}
