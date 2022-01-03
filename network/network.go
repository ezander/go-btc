package network

import (
	"net"
	"time"
)

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
	out = MarshalUint32(out, packet.Magic)
	out = MarshalFixedStr(out, packet.Command, 12)

	payload := packet.Message.Marshal([]byte{})
	out = MarshalUint32(out, uint32(len(payload)))
	out = MarshalUint32(out, checksum(payload))
	out = MarshalBytes(out, payload)
	return out
}

func UnmarshalPacket(data []byte) (*Packet, []byte) {
	origData := data
	if len(data) < 4+12+4+4 {
		return nil, origData
	}

	var packet Packet
	packet.Magic, data = UnmarshalUint32(data)
	packet.Command, data = UnmarshalFixedStr(data, 12)

	length, data := UnmarshalUint32(data)
	expectedChecksum, data := UnmarshalUint32(data)

	if len(data) < int(length) {
		return nil, origData
	}
	payload, data := UnmarshalBytes(data, length)

	actualPayload := checksum(payload)
	if expectedChecksum != actualPayload {
		panic("Checksums don't match")
	}

	message, data := unmarshalMessage(packet.Command, payload)
	packet.Message = message

	return &packet, data
}

func GetTestAddr(n int) net.TCPAddr {
	// https://bitcoin.stackexchange.com/questions/49634/testnet-peers-list-with-ip-addresses
	// dig A testnet-seed.bitcoin.jonasschnelli.ch

	ips, _ := net.LookupIP("testnet-seed.bitcoin.jonasschnelli.ch")

	return net.TCPAddr{IP: ips[n], Port: 18333, Zone: ""}
}

func GetTestConn(n int) net.Conn {
	tcp := GetTestAddr(n)
	conn, err := net.DialTimeout("tcp", tcp.String(), time.Millisecond*2000)
	if err != nil {
		panic(err)
	}
	// conn.SetDeadline(time.Now().Add(time.Second * 2))
	// fmt.Println(conn, err)
	return conn
}

type client struct {
	conn   net.Conn
	buffer []byte
	ready  bool
}

func Client(netConn net.Conn) client {
	return client{
		conn:   netConn,
		buffer: []byte{},
	}
}
func (cl *client) Close() error {
	return cl.conn.Close()
}
func (cl *client) ReadPacket() Packet {
	// fmt.Println("ReadPacket...")
	readBuf := make([]byte, 2048)
	for {
		// fmt.Println("Reading...")
		n, err := cl.conn.Read(readBuf)
		// fmt.Println(n, err)
		if err != nil {
			panic(err)
		}
		cl.buffer = append(cl.buffer, readBuf[:n]...)
		packet, buffer := UnmarshalPacket(cl.buffer)
		if packet != nil {
			cl.buffer = buffer
			return *packet
		}
	}
}
func (cl *client) SendPacket(packet Packet) {
	out := MarshalPacket(nil, packet)
	_, err := cl.conn.Write(out)
	if err != nil {
		panic(err)
	}
}
