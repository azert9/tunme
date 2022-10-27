package packet_link

import "io"

type streamManager struct {
	newStreamsChan chan struct{}
}

func newStreamManager() *streamManager {
	return &streamManager{
		newStreamsChan: make(chan struct{}),
	}
}

func (mgr *streamManager) HandlePacket(packet []byte) {
	mgr.newStreamsChan <- struct{}{}
}

func (mgr *streamManager) Accept() (io.ReadWriteCloser, error) {

	_ = <-mgr.newStreamsChan

	return nil, nil
}
