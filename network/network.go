package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func AsJSON(object interface{}) string {
	json, _ := json.MarshalIndent(object, "", "\t")
	return string(json)
}

const MAGIC_main uint32 = 0xD9B4BEF9
const MAGIC_testnet uint32 = 0xDAB5BFFA
const MAGIC_testnet3 uint32 = 0x0709110B
const MAGIC_signet uint32 = 0x40CF030A
const MAGIC_namecoin uint32 = 0xFEB4BEF9

// ======================================================================

type Packet struct {
	Magic   uint32
	Command string
	Message Message
}

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

func cleanString(in []byte) string {
	out := make([]byte, len(in))
	for i, c := range in {
		out[i] = c
		if c < ' ' || c > 0x7F {
			out[i] = '.'
		}
	}
	return string(out)
}
func hexPrinter(in []byte) {
	l := 16
	for start := 0; start < len(in); start += l {
		fmt.Printf("0x%04X  % x  %s\n", start, in[start:start+l], cleanString(in[start:start+l]))
	}
}

func UnmarshalPacket(data []byte, expectedMagic uint32) (*Packet, []byte) {
	origData := data
	if len(data) < 4+12+4+4 {
		return nil, origData
	}

	var packet Packet
	packet.Magic, data = UnmarshalUint32(data)
	if expectedMagic != 0 && packet.Magic != expectedMagic {
		fmt.Printf("Warning: Magic number mismatch (%x!=%x)'\n", expectedMagic, packet.Magic)
		hexPrinter(origData)
	}

	packet.Command, data = UnmarshalFixedStr(data, 12)

	length, data := UnmarshalUint32(data)
	expectedChecksum, data := UnmarshalUint32(data)

	if len(data) < int(length) {
		return nil, origData
	}
	payload, data := UnmarshalBytes(data, length)

	actualChecksum := checksum(payload)
	if expectedChecksum != actualChecksum {
		fmt.Printf("Warning: Checksums don't match (%x!=%x) in %v\n", expectedChecksum, actualChecksum, packet.Command)
	}

	message, payload := unmarshalMessage(packet.Command, payload)
	packet.Message = message
	if len(payload) > 0 {
		fmt.Printf("Warning: payload in message '%v' not fully used.", packet.Command)
	}

	return &packet, data
}

// ======================================================================

type client struct {
	conn   net.Conn
	magic  uint32
	buffer []byte
}

func Client(netConn net.Conn, magic uint32) client {
	return client{
		conn:   netConn,
		magic:  magic,
		buffer: []byte{},
	}
}

func (cl *client) Close() error {
	return cl.conn.Close()
}

func (cl *client) ReadPacket() Packet {
	readBuf := make([]byte, 2048)
	for {
		// See whether we have a complete packet in our buffer
		packet, buffer := UnmarshalPacket(cl.buffer, cl.magic)
		if packet != nil {
			cl.buffer = buffer
			return *packet
		}

		// Otherwise keep on reading from the tcp stream
		n, err := cl.conn.Read(readBuf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d bytes...\n", n)
		cl.buffer = append(cl.buffer, readBuf[:n]...)
	}
}

func (cl *client) SendPacket(packet Packet) {
	out := MarshalPacket(nil, packet)
	_, err := cl.conn.Write(out)
	if err != nil {
		panic(err)
	}
}

func (cl *client) SendMessage(message Message) {
	command := message.GetCommandString()
	packet := CreatePacket(cl.magic, command, message)
	fmt.Println("Sending: ", AsJSON(packet))
	cl.SendPacket(packet)
}

func (cl *client) ReceiveMessage() (*Message, string) {
	packet := cl.ReadPacket()
	if cl.magic != packet.Magic {
		fmt.Printf("Warning: magic bytes did not match: %x != %x\n", cl.magic, packet.Magic)
	}
	return &packet.Message, packet.Command
}

// ================================================================================================
func GetPeerAddress(seed string, port, n int) net.TCPAddr {
	ips, _ := net.LookupIP(seed)
	return net.TCPAddr{IP: ips[n], Port: port, Zone: ""}
}

func GetConnection(seed string, port int, n int) net.Conn {
	tcp := GetPeerAddress(seed, port, n)
	conn, err := net.DialTimeout("tcp", tcp.String(), time.Millisecond*2000)
	if err != nil {
		panic(err)
	}
	// fmt.Println(conn, err)
	return conn
}

// ==============================================================================================

func TestClient(n int) client {
	// https://bitcoin.stackexchange.com/questions/49634/testnet-peers-list-with-ip-addresses
	// dig A testnet-seed.bitcoin.jonasschnelli.ch

	seed := "testnet-seed.bitcoin.jonasschnelli.ch"
	port := 18333

	conn := GetConnection(seed, port, n)
	return Client(conn, MAGIC_testnet3)
}
