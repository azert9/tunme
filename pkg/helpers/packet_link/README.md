## Packet Format

Packet:

* `packet_type: Uint8`
* ...

PacketPacket:

* `packet_type: Uint8 = 1`
* `payload: Byte[]`

StreamDataPacket:

* `packet_type: Uint8 = 2`
* `stream_id: Uint32`
* `stream_offset: Uint64`
* `payload: Byte[]`

StreamAckPacket:

* `packet_type: Uint8 = 3`
* `stream_id: Uint32`
* `stream_offset: Uint64`

## TODO

* Use multiple streams for performances.
