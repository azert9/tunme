package packet_link

import (
	"io"
	"log"
	"sync"
	"time"
	"tunme/pkg/link"
)

// TODO: handle closed streams (remove from the map)

type streamManager struct {
	_mutex          sync.Mutex // TODO: reduce lock contention with a concurrent map
	_newStreamsChan chan *stream
	_streams        map[streamId]*stream
	_nextStreamId   streamId
	_packetSender   link.PacketSender
}

func newStreamManager(isServer bool, packetSender link.PacketSender) *streamManager {

	var firstStreamId streamId
	if isServer {
		firstStreamId = 0x80000000
	}

	return &streamManager{
		_streams:        make(map[streamId]*stream),
		_newStreamsChan: make(chan *stream, 10),
		_nextStreamId:   firstStreamId,
		_packetSender:   packetSender,
	}
}

func (mgr *streamManager) handleReceivedDataPacket(packet dataPacket) {

	mgr._mutex.Lock()
	defer mgr._mutex.Unlock()

	stream, found := mgr._streams[packet.getStreamId()]
	if !found {
		stream = newStream(packet.getStreamId(), mgr._packetSender)
		mgr._streams[stream.id] = stream
		mgr._newStreamsChan <- stream // TODO: handle full channel
	}

	// Even if the payload is empty, this will trigger an ACK.
	stream.handleReceivedDataPacket(packet)
}

func (mgr *streamManager) handleReceivedAckPacket(packet ackPacket) {

	mgr._mutex.Lock()
	defer mgr._mutex.Unlock()

	stream, found := mgr._streams[packet.getStreamId()]
	if !found {
		log.Printf("received ack packet for unknown stream")
		return
	}

	stream.handleReceivedAckPacket(packet)
}

func (mgr *streamManager) _newStream() *stream {

	mgr._mutex.Lock()
	defer mgr._mutex.Unlock()

	// TODO: detect streamId overflow
	streamId := mgr._nextStreamId
	mgr._nextStreamId++

	stream := newStream(streamId, mgr._packetSender)

	mgr._streams[streamId] = stream

	return stream
}

func (mgr *streamManager) openStream() *stream {

	stream := mgr._newStream()

	packet := newDataPacket(0)
	packet.setStreamId(stream.id)

	for {

		if err := mgr._packetSender.SendPacket(packet.getBytes(0)); err != nil {
			panic(err) // TODO
		}

		select {
		case _ = <-stream.firstAckChan:
			return stream
		case _ = <-time.After(4 * time.Second): // TODO: configure
		}
	}
}

func (mgr *streamManager) accept() (io.ReadWriteCloser, error) {

	stream := <-mgr._newStreamsChan

	return stream, nil
}
