package circular_buffer

import (
	"errors"
	"io"
	"sync"
	"tunme/internal/utils"
)

var ErrNotEnoughSpaceInBuffer = errors.New("not enough space in buffer")

type CircularBuffer struct {
	_mutex         sync.Mutex
	_writeCond     *sync.Cond
	_readCond      *sync.Cond
	_buff          []byte
	_off           int
	_len           int
	_closed        bool
	_blockingWrite bool
}

func NewCircularBuffer(capacity int) *CircularBuffer {

	buff := &CircularBuffer{
		_buff: make([]byte, capacity),
	}

	buff._writeCond = sync.NewCond(&buff._mutex)
	buff._readCond = sync.NewCond(&buff._mutex)

	return buff
}

func (buff *CircularBuffer) Len() int {
	// TODO: ensure safe concurrent access
	return buff._len
}

func (buff *CircularBuffer) Capacity() int {
	return len(buff._buff)
}

// SetBlockingWrite should be called by the writer (not concurrently with Write).
func (buff *CircularBuffer) SetBlockingWrite(blocking bool) {
	buff._blockingWrite = blocking
}

// Close always return a nil error.
func (buff *CircularBuffer) Close() error {

	buff._mutex.Lock()
	defer buff._mutex.Unlock()

	buff._closed = true

	buff._writeCond.Broadcast()

	return nil
}

// _writeSome will write as much data as possible, regarding the available space.
// The caller is responsible for holding the mutex, broadcasting conditions, and checking for closed flag.
func (buff *CircularBuffer) _writeSome(data []byte) int {

	availableSpace := len(buff._buff) - buff._len

	if len(data) > availableSpace {
		data = data[:availableSpace]
	}

	dataEnd := (buff._off + buff._len) % len(buff._buff)

	availableSpaceBeforeWrap := len(buff._buff) - dataEnd

	if len(data) > availableSpaceBeforeWrap {
		// wrapping
		copy(buff._buff[dataEnd:], data[:availableSpaceBeforeWrap])
		copy(buff._buff, data[availableSpaceBeforeWrap:])
	} else {
		copy(buff._buff[dataEnd:], data)
	}

	buff._len += len(data)

	return len(data)
}

func (buff *CircularBuffer) Write(data []byte) (int, error) {

	buff._mutex.Lock()
	defer buff._mutex.Unlock()

	offset := 0

	for offset < len(data) {

		// The closed flag might have been set while waiting for the condition.
		if buff._closed {
			return 0, io.EOF
		}

		n := buff._writeSome(data[offset:])

		if n == 0 {
			if buff._blockingWrite {
				buff._readCond.Wait()
			} else {
				return offset, ErrNotEnoughSpaceInBuffer
			}
		} else {
			offset += n
			buff._writeCond.Broadcast()
		}
	}

	if buff._closed {
		return len(data), io.EOF
	} else {
		return len(data), nil
	}
}

func (buff *CircularBuffer) read(out []byte, consume bool) (int, error) {

	if len(out) == 0 {
		return 0, nil
	}

	buff._mutex.Lock()
	defer buff._mutex.Unlock()

	for buff._len == 0 && !buff._closed {
		buff._writeCond.Wait()
	}

	if buff._len == 0 && buff._closed {
		return 0, io.EOF
	}

	rdLen := utils.Min(buff._len, len(out))

	availableDataBeforeWrap := len(buff._buff) - buff._off

	if rdLen > availableDataBeforeWrap {
		copy(out[:availableDataBeforeWrap], buff._buff[buff._off:])
		copy(out[availableDataBeforeWrap:], buff._buff[:])
	} else {
		copy(out, buff._buff[buff._off:])
	}

	if consume {
		buff._off = (buff._off + rdLen) % len(buff._buff)
		buff._len -= rdLen
		buff._readCond.Broadcast()
	}

	return rdLen, nil
}

func (buff *CircularBuffer) Read(out []byte) (int, error) {
	return buff.read(out, true)
}

// Peek behaves like Read, but does not consume data.
func (buff *CircularBuffer) Peek(out []byte) (int, error) {
	return buff.read(out, false)
}
