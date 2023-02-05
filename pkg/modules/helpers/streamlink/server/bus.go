package server

import "io"

type bus struct {
	acceptedStreamsChan   chan io.ReadWriteCloser
	receivedPacketsChan   chan []byte
	outControlPacketsChan chan []byte
	callBackStreamsChan   chan io.ReadWriteCloser
}

func (c *bus) closeAll() {
	close(c.acceptedStreamsChan)
	close(c.receivedPacketsChan)
	close(c.outControlPacketsChan)
	close(c.callBackStreamsChan)
}
