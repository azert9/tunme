## Packet Format

Packet:

* `packet_type: Uint8`
* ...

PacketPacket:

* `packet_type: Uint8 = 0`
* `payload: Byte[]`

StreamPacket:

* `packet_type: Uint8 = 1`
* `stream_id: Uint32`
* `stream_offset: Uint64`
* `payload: Byte[]`
