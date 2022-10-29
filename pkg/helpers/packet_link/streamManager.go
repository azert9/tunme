package packet_link

import (
	"io"
	"sync"
)

type streamManager struct {
	mutex          sync.Mutex
	newStreamsChan chan struct{}
	streams        map[streamId]*stream
}

func newStreamManager() *streamManager {
	return &streamManager{
		newStreamsChan: make(chan struct{}),
	}
}

func (mgr *streamManager) HandlePacket(packet []byte) {

	//if len(dataPacket) < 8 {
	//	log.Printf("dropping invalid dataPacket")
	//	return
	//}
	//streamId := _getUnsafe[streamId](dataPacket, 4)
	//streamOffset := _getUnsafe[uint64](dataPacket, 8)

	//// TODO: reduce lock congestion

	//mgr.mutex.Lock()
	//defer mgr.mutex.Unlock()

	//streamReceiver, found := mgr.streams[streamId]
	//if !found {

	//	if streamOffset != 0 {
	//		log.Printf("dropping invalid dataPacket")
	//		return
	//	}

	//	// TODO: create the streamReceiver

	//	mgr.newStreamsChan <- struct{}{}
	//}

	//// TODO: truncate the part of the dataPacket which

	//_ = streamReceiver

	// TODO: write the received data
	// TODO: send an ACK (which must contain information about the available window size)
}

func (mgr *streamManager) Accept() (io.ReadWriteCloser, error) {

	_ = <-mgr.newStreamsChan

	return nil, nil
}
