package network

type Packet struct {
	Magic   uint32
	Command string
	Message Message
}

const MAGIC_main = 0xD9B4BEF9
const MAGIC_testnet = 0xDAB5BFFA
const MAGIC_testnet3 = 0x0709110B
const MAGIC_signet = 0x40CF030A
const MAGIC_namecoin = 0xFEB4BEF9

func CreatePacket(magic uint32, command string, msg Message) Packet {
	message := Packet{magic, command, msg}
	return message
}

func MarshalPacket(out []byte, packet Packet) []byte {
	out = marshalUint32(out, packet.Magic)
	out = MarshalFixedStr(out, packet.Command, 12)

	payload := packet.Message.Marshal([]byte{})
	out = marshalUint32(out, uint32(len(payload)))
	out = marshalUint32(out, checksum(payload))
	out = MarshalBytes(out, payload)
	return out
}

func UnmarshalPacket(data []byte) (Packet, []byte) {
	var packet Packet
	packet.Magic, data = unmarshalUint32(data)
	packet.Command, data = UnmarshalFixedStr(data, 12)

	length, data := unmarshalUint32(data)
	expectedChecksum, data := unmarshalUint32(data)
	payload, data := UnmarshalBytes(data, length)

	actualPayload := checksum(payload)
	if expectedChecksum != actualPayload {
		panic("Checksums don't match")
	}

	message, data := unmarshalMessage(packet.Command, payload)
	packet.Message = message

	return packet, data
}
