package network

import (
	"fmt"
	"net"
	"time"
)

type Message interface {
	Marshaller
	Unmarshaller
	GetCommandString() string
}

func unmarshalMessage(command string, data []byte) (Message, []byte) {
	var msg Message
	switch command {
	case "version":
		msg = new(VersionMessage)
	case "verack":
		msg = new(VerAckMessage)
	case "reject":
		msg = new(RejectMessage)
	case "ping":
		msg = new(PingMessage)
	case "pong":
		msg = new(PongMessage)
	default:
		panic(fmt.Sprintf("Unknown command to unmarshal: '%s'", command))
	}
	if command != msg.GetCommandString() {
		panic("Internal error (command string mismatch)")
	}
	data = msg.Unmarshal(data)
	return msg, data
}

// ========================================================================

type VersionMessage struct {
	Version      uint32    // Identifies protocol version being used by the node
	Services     uint64    // bitfield of features to be enabled for this connection
	Timestamp    time.Time // standard UNIX timestamp in seconds
	ReceiverAddr NetAddr   // The network address of the node receiving this message
	// Fields below require version ≥ 106
	FromAddr    NetAddr // Field can be ignored. This used to be the network address of the node emitting this message, but most P2P implementations send 26 dummy bytes. The "services" field of the address would also be redundant with the second field of the version message.
	Nonce       uint64  // 	Node random nonce, randomly generated every time a version packet is sent. This nonce is used to detect connections to self.
	UserAgent   string  //	User Agent (0x00 if string is 0 bytes long)
	StartHeight uint32  // The last block received by the emitting node
	// Fields below require version ≥ 70001
	Relay bool // 	Whether the remote peer should announce relayed transactions or not, see BIP 0037
}

// Bitfield constants for services

const NODE_NETWORK = 1            // 	This node can be asked for full blocks instead of just headers.
const NODE_GETUTXO = 2            // 	See BIP 0064
const NODE_BLOOM = 4              // 	See BIP 0111
const NODE_WITNESS = 8            // 	See BIP 0144
const NODE_XTHIN = 16             // 	Never formally proposed (as a BIP), and discontinued. Was historically sporadically seen on the network.
const NODE_COMPACT_FILTERS = 64   // 	See BIP 0157
const NODE_NETWORK_LIMITED = 1024 // 	See BIP 0159

func NewVersionMessage() *VersionMessage {

	msg := VersionMessage{
		Version:      31800,
		Services:     NODE_NETWORK,
		Timestamp:    time.Now().Truncate(time.Second),
		ReceiverAddr: NetAddr{NODE_NETWORK, net.IPv4(127, 0, 0, 1), 8333},
		FromAddr:     NetAddr{NODE_NETWORK, net.IPv4(127, 0, 0, 1), 8333},
		Nonce:        3141526,
		UserAgent:    "Foobar client v0.1",
		StartHeight:  1,
		Relay:        false,
	}
	return &msg
}

func (msg VersionMessage) Marshal(out []byte) []byte {
	out = MarshalUint32(out, msg.Version)
	out = MarshalUint64(out, msg.Services)
	out = MarshalTimestamp(out, msg.Timestamp)
	out = MarshalNetAddr(out, msg.ReceiverAddr)
	if msg.Version >= 106 {
		out = MarshalNetAddr(out, msg.FromAddr)
		out = MarshalUint64(out, msg.Nonce)
		out = MarshalVarStr(out, msg.UserAgent)
		out = MarshalUint32(out, msg.StartHeight)
	}
	if msg.Version >= 70001 {
		out = MarshalBool(out, msg.Relay)
	}
	return out
}

func (msg *VersionMessage) Unmarshal(data []byte) []byte {
	msg.Version, data = UnmarshalUint32(data)
	msg.Services, data = UnmarshalUint64(data)
	msg.Timestamp, data = UnmarshalTimestamp(data)
	msg.ReceiverAddr, data = UnmarshalNetAddr(data)
	if msg.Version >= 106 {
		msg.FromAddr, data = UnmarshalNetAddr(data)
		msg.Nonce, data = UnmarshalUint64(data)
		msg.UserAgent, data = UnmarshalVarStr(data)
		msg.StartHeight, data = UnmarshalUint32(data)
	}
	if msg.Version >= 70001 {
		msg.Relay, data = UnmarshalBool(data)
	}
	return data
}

func (msg VersionMessage) GetCommandString() string {
	return "version"
}

// ========================================================================
type VerAckMessage struct {
}

func (msg VerAckMessage) Marshal(out []byte) []byte {
	return out
}

func (msg *VerAckMessage) Unmarshal(data []byte) []byte {
	return data
}

func (msg VerAckMessage) GetCommandString() string {
	return "verack"
}

// ========================================================================

const REJECT_MALFORMED uint8 = 0x01
const REJECT_INVALID uint8 = 0x10
const REJECT_OBSOLETE uint8 = 0x11
const REJECT_DUPLICATE uint8 = 0x12
const REJECT_NONSTANDARD uint8 = 0x40
const REJECT_DUST uint8 = 0x41
const REJECT_INSUFFICIENTFEE uint8 = 0x42
const REJECT_CHECKPOINT uint8 = 0x43

type RejectMessage struct {
	Message string //	var_str - type of message rejected
	CCode   uint8  // char - code relating to rejected message
	Reason  string // var_str - text version of reason for rejection
	Data    []byte // char - Optional extra data provided by some errors. Currently, all errors which provide this field fill it with the TXID or block header hash of the object being rejected, so the field is 32 bytes.
}

func (msg RejectMessage) Marshal(out []byte) []byte {
	out = MarshalVarStr(out, msg.Message)
	out = MarshalUint8(out, msg.CCode)
	out = MarshalVarStr(out, msg.Reason)
	out = MarshalBytes(out, msg.Data[:])
	return out
}

func (msg *RejectMessage) Unmarshal(data []byte) []byte {
	msg.Message, data = UnmarshalVarStr(data)
	msg.CCode, data = UnmarshalUint8(data)
	msg.Reason, data = UnmarshalVarStr(data)
	msg.Data, data = UnmarshalBytes(data, uint32(len(data)))
	return data
}

func (msg RejectMessage) GetCommandString() string {
	return "reject"
}

// ========================================================================
type PingMessage struct {
	Nonce uint64
}

func (msg PingMessage) Marshal(out []byte) []byte {
	out = MarshalUint64(out, msg.Nonce)
	return out
}

func (msg *PingMessage) Unmarshal(data []byte) []byte {
	if len(data) >= 8 {
		msg.Nonce, data = UnmarshalUint64(data)
	}
	return data
}

func (msg PingMessage) GetCommandString() string {
	return "ping"
}

// ========================================================================
type PongMessage struct {
	Nonce uint64
}

func (msg PongMessage) Marshal(out []byte) []byte {
	out = MarshalUint64(out, msg.Nonce)
	return out
}

func (msg *PongMessage) Unmarshal(data []byte) []byte {
	msg.Nonce, data = UnmarshalUint64(data)
	return data
}

func (msg PongMessage) GetCommandString() string {
	return "pong"
}

// type AddrMessage struct {
// }
