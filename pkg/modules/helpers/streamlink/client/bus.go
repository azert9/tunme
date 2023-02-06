package client

// TODO: use something similar for the server

type bus struct {
	_closeChan  chan struct{}
	_acceptChan chan struct{}
}

func newBus() *bus {
	return &bus{
		_closeChan:  make(chan struct{}),
		_acceptChan: make(chan struct{}),
	}
}

func (bus *bus) close() {
	close(bus._closeChan)
}

func (bus *bus) sendAcceptNonBlocking() bool {

	select {
	case bus._acceptChan <- struct{}{}:
		return true
	case <-bus._closeChan:
		return false
	default:
		return false
	}
}

func (bus *bus) receiveAccept() bool {

	select {
	case <-bus._acceptChan:
		return true
	case <-bus._closeChan:
		return false
	}
}
