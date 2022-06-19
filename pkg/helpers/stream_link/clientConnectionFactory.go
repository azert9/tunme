package stream_link

import "net"

type _clientConnectionFactory struct {
	Dialer StreamDialer
}

func newClientConnectionFactory(dialer StreamDialer) connectionFactory {
	return &_clientConnectionFactory{
		Dialer: dialer,
	}
}

func (c *_clientConnectionFactory) MakeConnection(channel int) (net.Conn, error) {

	conn, err := c.Dialer.Dial()
	if err != nil {
		return nil, err
	}

	var buff [1]byte
	buff[0] = byte(channel)

	if _, err := conn.Write(buff[:]); err != nil {
		return nil, err
	}

	return conn, nil
}
