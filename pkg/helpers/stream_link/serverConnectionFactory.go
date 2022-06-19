package stream_link

import (
	"io"
	"net"
)

// TODO: handle closing

type _serverConnectionFactory struct {
	Listener net.Listener
	Channels []chan net.Conn
}

func newServerConnectionFactory(listener net.Listener, channelCount int) connectionFactory {

	f := &_serverConnectionFactory{
		Listener: listener,
	}

	f.Channels = make([]chan net.Conn, channelCount)
	for i := range f.Channels {
		f.Channels[i] = make(chan net.Conn)
	}

	go f._acceptLoop()

	return f
}

func (f *_serverConnectionFactory) _acceptLoop() {

	for {

		conn, err := f.Listener.Accept()
		if err != nil {
			// TODO: Warn or something? Or fail, but then retry should be implemented earlier.
			continue
		}

		var buff [1]byte
		if _, err := io.ReadFull(conn, buff[:]); err != nil {
			conn.Close() // TODO: warn
			continue
		}

		channel := int(buff[0])
		if channel >= len(f.Channels) {
			// TODO: warn
			continue
		}

		f.Channels[channel] <- conn
	}
}

func (f *_serverConnectionFactory) MakeConnection(channel int) (net.Conn, error) {
	conn := <-f.Channels[channel]
	return conn, nil
}
