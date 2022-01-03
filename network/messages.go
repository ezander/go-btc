package network

import (
	"fmt"
	"net"
	"time"
)

type Message interface {
	Marshaller
	Unmarshaller
}

func unmarshalMessage(command string, data []byte) (Message, []byte) {
	var msg Message
	switch command {
	case "version":
		msg = new(VersionMessage)
	default:
		panic(fmt.Sprintf("Unknown command to unmarshal: '%s'", command))
	}
	data = msg.Unmarshal(data)
	return msg, data
}

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

func NewVersionMessage() VersionMessage {

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
	return msg
}

func (v VersionMessage) Marshal(out []byte) []byte {
	out = marshalUint32(out, v.Version)
	out = marshalUint64(out, v.Services)
	out = MarshalTimestamp(out, v.Timestamp)
	out = MarshalNetAddr(out, v.ReceiverAddr)
	if v.Version >= 106 {
		out = MarshalNetAddr(out, v.FromAddr)
		out = MarshalNonce(out, v.Nonce)
		out = MarshalVarStr(out, v.UserAgent)
		out = marshalUint32(out, v.StartHeight)
	}
	if v.Version >= 70001 {
		out = MarshalBool(out, v.Relay)
	}
	return out
}

func (msg *VersionMessage) Unmarshal(out []byte) []byte {
	msg.Version, out = unmarshalUint32(out)
	msg.Services, out = unmarshalUint64(out)
	msg.Timestamp, out = UnmarshalTimestamp(out)
	msg.ReceiverAddr, out = UnmarshalNetAddr(out)
	if msg.Version >= 106 {
		msg.FromAddr, out = UnmarshalNetAddr(out)
		msg.Nonce, out = UnmarshalNonce(out)
		msg.UserAgent, out = UnmarshalVarStr(out)
		msg.StartHeight, out = unmarshalUint32(out)
	}
	if msg.Version >= 70001 {
		msg.Relay, out = UnmarshalBool(out)
	}
	return out
}

func UnmarshalVersionMessage(out []byte) (VersionMessage, []byte) {
	msg := VersionMessage{}
	out = msg.Unmarshal(out)
	return msg, out
}

// type VerackMessage struct {
// }

// type AddrMessage struct {
// }
