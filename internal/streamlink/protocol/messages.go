package protocol

type StreamType byte

const (
	StreamTypeControl StreamType = 0
	StreamTypeConnect StreamType = 1
	// StreamTypeCallBack is for when the client opens a connection at the request of the server.
	StreamTypeCallBack StreamType = 2
)

// ClientHello is sent by the client upon opening a connection.
type ClientHello struct {
	StreamType StreamType
}

type ControlPacketType byte

const (
	ControlPacketTypeData          ControlPacketType = 0
	ControlPacketTypeStreamRequest ControlPacketType = 1
)
