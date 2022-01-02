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
	version   uint32    // Identifies protocol version being used by the node
	services  uint64    // bitfield of features to be enabled for this connection
	timestamp time.Time // standard UNIX timestamp in seconds
	addr_recv NetAddr   // The network address of the node receiving this message
	//Fields below require version ≥ 106
	addr_from    NetAddr // Field can be ignored. This used to be the network address of the node emitting this message, but most P2P implementations send 26 dummy bytes. The "services" field of the address would also be redundant with the second field of the version message.
	nonce        uint64  // 	Node random nonce, randomly generated every time a version packet is sent. This nonce is used to detect connections to self.
	user_agent   string  //	User Agent (0x00 if string is 0 bytes long)
	start_height uint32  // The last block received by the emitting node
	// Fields below require version ≥ 70001
	relay bool // 	Whether the remote peer should announce relayed transactions or not, see BIP 0037
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
		version:      31800,
		services:     NODE_NETWORK,
		timestamp:    time.Now().Truncate(time.Second),
		addr_recv:    NetAddr{NODE_NETWORK, net.IPv4(127, 0, 0, 1), 8333},
		addr_from:    NetAddr{NODE_NETWORK, net.IPv4(127, 0, 0, 1), 8333},
		nonce:        3141526,
		user_agent:   "Foobar client v0.1",
		start_height: 1,
		relay:        false,
	}
	return msg
}

func (v VersionMessage) Marshal(out []byte) []byte {
	out = marshalUint32(out, v.version)
	out = marshalUint64(out, v.services)
	out = MarshalTimestamp(out, v.timestamp)
	out = MarshalNetAddr(out, v.addr_recv)
	if v.version >= 106 {
		out = MarshalNetAddr(out, v.addr_from)
		out = MarshalNonce(out, v.nonce)
		out = MarshalVarStr(out, v.user_agent)
		out = marshalUint32(out, v.start_height)
	}
	if v.version >= 70001 {
		out = MarshalBool(out, v.relay)
	}
	return out
}

func (msg *VersionMessage) Unmarshal(out []byte) []byte {
	msg.version, out = unmarshalUint32(out)
	msg.services, out = unmarshalUint64(out)
	msg.timestamp, out = UnmarshalTimestamp(out)
	msg.addr_recv, out = UnmarshalNetAddr(out)
	if msg.version >= 106 {
		msg.addr_from, out = UnmarshalNetAddr(out)
		msg.nonce, out = UnmarshalNonce(out)
		msg.user_agent, out = UnmarshalVarStr(out)
		msg.start_height, out = unmarshalUint32(out)
	}
	if msg.version >= 70001 {
		msg.relay, out = UnmarshalBool(out)
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
