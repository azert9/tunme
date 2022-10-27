package circular_buffer

import (
	"errors"
	"sync"
	"tunme/internal/utils"
)

var ErrNotEnoughSpaceInBuffer = errors.New("not enough space in buffer")

type CircularBuffer struct {
	_mutex sync.Mutex
	_buff  []byte
	_off   int
	_len   int
}

func NewCircularBuffer(capacity int) *CircularBuffer {
	return &CircularBuffer{
		_buff: make([]byte, capacity),
	}
}

func (buff *CircularBuffer) Capacity() int {
	return len(buff._buff)
}

func (buff *CircularBuffer) Write(data []byte) (int, error) {

	buff._mutex.Lock()
	defer buff._mutex.Unlock()

	// TODO: detect wrapped writes

	dataEnd := (buff._off + buff._len) % len(buff._buff)

	availableSpace := len(buff._buff) - buff._len

	if len(data) > availableSpace {
		return 0, ErrNotEnoughSpaceInBuffer
	}

	availableSpaceBeforeWrap := len(buff._buff) - dataEnd

	if len(data) > availableSpaceBeforeWrap {
		// wrapping
		copy(buff._buff[dataEnd:], data[:availableSpaceBeforeWrap])
		copy(buff._buff, data[availableSpaceBeforeWrap:])
	} else {
		copy(buff._buff[dataEnd:], data)
	}

	buff._len += len(data)

	return len(data), nil
}

func (buff *CircularBuffer) Read(out []byte) (int, error) {

	buff._mutex.Lock()
	defer buff._mutex.Unlock()

	rdLen := utils.Min(buff._len, len(out))

	availableDataBeforeWrap := len(buff._buff) - buff._off

	if rdLen > availableDataBeforeWrap {
		copy(out[:availableDataBeforeWrap], buff._buff[buff._off:])
		copy(out[availableDataBeforeWrap:], buff._buff[:])
	} else {
		copy(out, buff._buff[buff._off:])
	}

	buff._off = (buff._off + rdLen) % len(buff._buff)
	buff._len -= rdLen

	return rdLen, nil
}
