package conngc

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type ConnGarbageCollector interface {
	CloseAll()
	OpenConn(f func() (io.ReadWriteCloser, error)) (io.ReadWriteCloser, error)
}

type impl struct {
	mutex       sync.Mutex
	connections map[uint64]io.ReadWriteCloser
	nextId      uint64
}

func New() ConnGarbageCollector {
	return &impl{
		connections: map[uint64]io.ReadWriteCloser{},
	}
}

// CloseAll is safe to call multiple times.
func (ct *impl) CloseAll() {

	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if ct.connections == nil {
		return
	}

	for _, conn := range ct.connections {
		conn.Close()
	}

	ct.connections = nil
}

func (ct *impl) OpenConn(f func() (io.ReadWriteCloser, error)) (io.ReadWriteCloser, error) {

	conn, err := f()
	if err != nil {
		return nil, err
	}

	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if ct.connections == nil {
		if conn.Close(); err != nil {
			log.Print(err)
		}
		return nil, fmt.Errorf("tunnel closed") // TODO: proper error
	}

	id := ct.nextId
	ct.nextId++

	ct.connections[id] = conn

	return &trackedConn{
		id:      id,
		tracker: ct,
		conn:    conn,
	}, nil
}

func (ct *impl) remove(id uint64) {

	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	delete(ct.connections, id)
}

type trackedConn struct {
	id      uint64
	tracker *impl
	conn    io.ReadWriteCloser
}

func (c *trackedConn) Read(buff []byte) (int, error) {
	return c.conn.Read(buff)
}

func (c *trackedConn) Write(buff []byte) (int, error) {
	return c.conn.Write(buff)
}

func (c *trackedConn) Close() error {
	c.tracker.remove(c.id)
	return c.conn.Close()
}
