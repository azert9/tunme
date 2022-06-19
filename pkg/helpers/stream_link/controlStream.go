package stream_link

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sync"
	"time"
)

type controlStream interface {
	Send(packet []byte) error
	Receive() ([]byte, error)
	SendStreamRequest() error
	WaitStreamRequest()
	//Close() error // TODO
}

type _controlStream struct {
	ReceiveChan       chan _received
	StreamRequestChan chan any
	WriterMutex       sync.Mutex
	Writer            io.Writer
	WriterReadyCond   *sync.Cond // condition: Writer is non-nil
}

type _received struct {
	Packet []byte
	Err    error
}

func newControlStream(connectionFactory connectionFactory) controlStream {

	stream := &_controlStream{
		ReceiveChan:       make(chan _received),
		StreamRequestChan: make(chan any),
	}
	stream.WriterReadyCond = sync.NewCond(&stream.WriterMutex)

	go func() {

		for {

			conn, err := connectionFactory.MakeConnection(0)
			if err != nil {
				// TODO: do not block by sending the error over the channel
				stream.ReceiveChan <- _received{
					Err: err,
				}
				time.Sleep(time.Second) // TODO: configure retry delay
				continue
			}

			stream.WriterMutex.Lock()
			stream.Writer = conn
			stream.WriterMutex.Unlock()
			stream.WriterReadyCond.Broadcast()

			err = _readPackets(conn, stream.ReceiveChan, stream.StreamRequestChan)
			// TODO: do not block by sending the error over the channel
			stream.ReceiveChan <- _received{
				Err: err,
			}
		}
	}()

	return stream
}

type _packetHeader struct {
	Len  int32
	Type uint8
}

const controlPacketTypeDatagram = 0
const controlPacketTypeStream = 1

func _readPackets(reader io.Reader, receiveChan chan<- _received, streamRequestChan chan<- any) error {

	for {
		var header _packetHeader
		headerBuff := make([]byte, binary.Size(header))
		if _, err := io.ReadFull(reader, headerBuff); err != nil {
			return err
		}
		if err := binary.Read(bytes.NewReader(headerBuff), binary.LittleEndian, &header); err != nil {
			return err
		}

		if header.Len < 0 {
			return fmt.Errorf("invalid packet header: length field is negative")
		}

		packet := make([]byte, header.Len)

		if _, err := io.ReadFull(reader, packet); err != nil {
			return err
		}

		switch header.Type {
		case controlPacketTypeDatagram:
			receiveChan <- _received{Packet: packet}
		case controlPacketTypeStream:
			streamRequestChan <- nil
		default:
			// TODO: warn
		}
	}
}

func (s *_controlStream) Send(packet []byte) error {

	s.WriterMutex.Lock()
	defer s.WriterMutex.Unlock()

	for s.Writer == nil {
		s.WriterReadyCond.Wait()
	}

	if len(packet) > math.MaxInt32 {
		return errPacketTooBig
	}
	header := _packetHeader{
		Len:  int32(len(packet)),
		Type: controlPacketTypeDatagram,
	}

	buff := bytes.NewBuffer(make([]byte, 0, binary.Size(header)+len(packet)))
	if err := binary.Write(buff, binary.LittleEndian, header); err != nil {
		return err
	}
	if _, err := buff.Write(packet); err != nil {
		return err
	}

	if _, err := s.Writer.Write(buff.Bytes()); err != nil {
		return err
	}

	return nil
}

// Receive may fail, in which case it can be re-called to try again.
func (s *_controlStream) Receive() ([]byte, error) {

	received := <-s.ReceiveChan

	if received.Err != nil {
		return nil, received.Err
	}

	return received.Packet, nil
}

func (s *_controlStream) SendStreamRequest() error {

	s.WriterMutex.Lock()
	defer s.WriterMutex.Unlock()

	for s.Writer == nil {
		s.WriterReadyCond.Wait()
	}

	header := _packetHeader{
		Type: controlPacketTypeStream,
	}

	buff := bytes.NewBuffer(make([]byte, 0, binary.Size(header)))
	if err := binary.Write(buff, binary.LittleEndian, header); err != nil {
		return err
	}

	if _, err := s.Writer.Write(buff.Bytes()); err != nil {
		return err
	}

	return nil
}

func (s *_controlStream) WaitStreamRequest() {
	<-s.StreamRequestChan
}
