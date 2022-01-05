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
	case "alert":
		msg = new(AlertMessage)
	case "addr":
		msg = new(AddrMessage)
	case "sendheaders":
		msg = new(SendHeadersMessage)
	case "getheaders":
		msg = new(GetHeadersMessage)
	case "getblocks":
		msg = new(GetBlocksMessage)
	case "inv":
		msg = new(InvMessage)
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
		Services:     0, //NODE_NETWORK,
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

// ========================================================================
type AlertMessage struct {
	Version    uint32    // int32_t 	Alert format version
	RelayUntil time.Time // int64_t The timestamp beyond which nodes should stop relaying this alert
	Expiration time.Time //	int64_t 	The timestamp beyond which this alert is no longer in effect and should be ignored
	ID         uint32    // int32_t 	A unique ID number for this alert
	Cancel     uint32    // int32_t 	All alerts with an ID number less than or equal to this number should be cancelled: deleted and not accepted in the future
	setCancel  []uint32  // set<int32_t> 	All alert IDs contained in this set should be cancelled as above
	MinVer     uint32    // int32_t 	This alert only applies to versions greater than or equal to this version. Other versions should still relay it.
	MaxVer     uint32    // int32_t 	This alert only applies to versions less than or equal to this version. Other versions should still relay it.
	setSubVer  []string  // set<string> 	If this set contains any elements, then only nodes that have their subVer contained in this set are affected by the alert. Other versions should still relay it.
	Priority   uint32    // int32_t 	Relative priority compared to other alerts
	Comment    string    // string 	A comment on the alert that is not displayed
	StatusBar  string    // string 	The alert message that is displayed to the user
	Reserved   string    // string 	Reserved
	Payload    []byte
	Signature  []byte
}

func (msg AlertMessage) Marshal(out []byte) []byte {
	panic("I ain't gonna send no alert message (see https://bitcoin.org/en/alert/2016-11-01-alert-retirement) ")
}

func (msg *AlertMessage) Unmarshal(data []byte) []byte {

	// Unmarshal payload and signature
	lenPayload, data := UnmarshalVarInt(data)
	msg.Payload, data = UnmarshalBytes(data, uint32(lenPayload))
	lenSignature, data := UnmarshalVarInt(data)
	msg.Signature, data = UnmarshalBytes(data, uint32(lenSignature))

	// Check signature
	// Public key
	//   04fc9702847840aaf195de8442ebecedf5b095cdbb9bc716bda9110971b28a49e0ead8564ff0db22209e0374782c093bb899692d524e9d6a6956e7c5ecbcd68284
	//   (hash) 1AGRxqDa5WjUKBwHB9XYEjmkv1ucoUUy1s
	// Public keys (in scripts) are given as
	// 		04 <x> <y> where x and y are 32 byte big-endian integers representing the coordinates of a point on the curve
	// or in compressed form given as
	//     <sign> <x> where <sign> is 0x02 if y is even and 0x03 if y is odd.
	key := PublicKeyFromString("04fc9702847840aaf195de8442ebecedf5b095cdbb9bc716bda9110971b28a49e0ead8564ff0db22209e0374782c093bb899692d524e9d6a6956e7c5ecbcd68284")
	hash := doubleHash(msg.Payload)
	valid := VerifyASN1(&key, hash[:], msg.Signature)
	if !valid {
		// we do nothing... does not matter anyway, as the alert
		// system has been retired and the private key is in the wild
		// (I just wanted to try whether I could validate the alert message
		// about retiring the alert system - didn't work though...)
		// And furthermore, it's send only to pre-70000 clients anyway
	}

	// Unmarshal fields from payload
	payload := msg.Payload
	msg.Version, payload = UnmarshalUint32(payload)
	msg.RelayUntil, payload = UnmarshalTimestamp(payload)
	msg.Expiration, payload = UnmarshalTimestamp(payload)
	msg.ID, payload = UnmarshalUint32(payload)
	msg.Cancel, payload = UnmarshalUint32(payload)
	nCancel, payload := UnmarshalVarInt(payload)
	msg.setCancel = make([]uint32, nCancel)
	for i := range msg.setCancel {
		msg.setCancel[i], payload = UnmarshalUint32(payload)
	}
	msg.MinVer, payload = UnmarshalUint32(payload)
	msg.MaxVer, payload = UnmarshalUint32(payload)
	nSubVer, payload := UnmarshalVarInt(payload)
	msg.setSubVer = make([]string, nSubVer)
	for i := range msg.setSubVer {
		msg.setSubVer[i], payload = UnmarshalVarStr(payload)
	}
	msg.Priority, payload = UnmarshalUint32(payload)
	msg.Comment, payload = UnmarshalVarStr(payload)
	msg.StatusBar, payload = UnmarshalVarStr(payload)
	msg.Reserved, payload = UnmarshalVarStr(payload)
	if len(payload) > 0 {
		fmt.Printf("Warning: Payload not fully parsed at end of alert message...")
	}

	return data
}

func (msg AlertMessage) GetCommandString() string {
	return "alert"
}

// ========================================================================
type AddrMessage struct {
	AddrList []TimeNetAddr
}

func (msg AddrMessage) Marshal(out []byte) []byte {
	out = MarshalVarInt(out, uint64(len(msg.AddrList)))
	for i := range msg.AddrList {
		out = MarshalTimeNetAddr(out, msg.AddrList[i])
	}
	return out
}

func (msg *AddrMessage) Unmarshal(data []byte) []byte {
	count, data := UnmarshalVarInt(data)
	msg.AddrList = make([]TimeNetAddr, count)
	for i := range msg.AddrList {
		msg.AddrList[i], data = UnmarshalTimeNetAddr(data)
	}
	return data
}

func (msg AddrMessage) GetCommandString() string {
	return "addr"
}

// ========================================================================

// sendheaders
//
// Request for Direct headers announcement.
//
// Upon receipt of this message, the node is be permitted, but not required, to announce new blocks by headers command (instead of inv command).
//
// This message is supported by the protocol version >= 70012 or Bitcoin Core version >= 0.12.0.
//
// See BIP 130 for more information.
//
// No additional data is transmitted with this message.
type SendHeadersMessage struct {
}

func (msg SendHeadersMessage) Marshal(out []byte) []byte {
	return out
}

func (msg *SendHeadersMessage) Unmarshal(data []byte) []byte {
	return data
}

func (msg SendHeadersMessage) GetCommandString() string {
	return "sendheaders"
}

// ========================================================================

// getheaders

// Return a headers packet containing the headers of blocks starting right after the last known hash in the block locator object, up to hash_stop or 2000 blocks, whichever comes first. To receive the next block headers, one needs to issue getheaders again with a new block locator object. Keep in mind that some clients may provide headers of blocks which are invalid if the block locator object contains a hash on the invalid branch.

// Payload:
// Field Size 	Description 	Data type 	Comments
// 4 	version 	uint32_t 	the protocol version
// 32+ 	block locator hashes 	char[32] 	block locator object; newest back to genesis block (dense to start, but then sparse)
// 32 	hash_stop 	char[32] 	hash of the last desired block header; set to zero to get as many blocks as possible (2000)

type GetHeadersMessage struct {
	Version        uint32
	BlockLocHashes []Hash
	StopHash       Hash
}

func (msg GetHeadersMessage) Marshal(out []byte) []byte {
	out = MarshalUint32(out, msg.Version)
	out = MarshalHashes(out, msg.BlockLocHashes)
	out = MarshalHash(out, msg.StopHash)
	return out
}

func (msg *GetHeadersMessage) Unmarshal(data []byte) []byte {
	msg.Version, data = UnmarshalUint32(data)
	msg.BlockLocHashes, data = UnmarshalHashes(data)
	msg.StopHash, data = UnmarshalHash(data)
	return data
}

func (msg GetHeadersMessage) GetCommandString() string {
	return "getheaders"
}

// ========================================================================

// getblocks

// Return an inv packet containing the list of blocks starting right after the last known hash in the block locator object, up to hash_stop or 500 blocks, whichever comes first.

// The locator hashes are processed by a node in the order as they appear in the message. If a block hash is found in the node's main chain, the list of its children is returned back via the inv message and the remaining locators are ignored, no matter if the requested limit was reached, or not.

// To receive the next blocks hashes, one needs to issue getblocks again with a new block locator object. Keep in mind that some clients may provide blocks which are invalid if the block locator object contains a hash on the invalid branch.

type GetBlocksMessage struct {
	// Payload:
	// Field Size 	Description 	Data type 	Comments
	// 4 	version 	uint32_t 	the protocol version
	// 1+ 	hash count 	var_int 	number of block locator hash entries
	// 32+ 	block locator hashes 	char[32] 	block locator object; newest back to genesis block (dense to start, but then sparse)
	// 32 	hash_stop 	char[32] 	hash of the last desired block; set to zero to get as many blocks as possible (500)
	Version        uint32
	BlockLocHashes []Hash
	StopHash       Hash
}

func (msg GetBlocksMessage) Marshal(out []byte) []byte {
	out = MarshalUint32(out, msg.Version)
	out = MarshalHashes(out, msg.BlockLocHashes)
	out = MarshalHash(out, msg.StopHash)
	return out
}

func (msg *GetBlocksMessage) Unmarshal(data []byte) []byte {
	msg.Version, data = UnmarshalUint32(data)
	msg.BlockLocHashes, data = UnmarshalHashes(data)
	msg.StopHash, data = UnmarshalHash(data)
	return data
}

func (msg GetBlocksMessage) GetCommandString() string {
	return "getblocks"
}

// ========================================================================

// Allows a node to advertise its knowledge of one or more objects. It can be
// received unsolicited, or in reply to getblocks.

// Payload (maximum 50,000 entries, which is just over 1.8 megabytes):
// Field Size 	Description 	Data type 	Comments
// 1+ 	count 	var_int 	Number of inventory entries
// 36x? 	inventory 	inv_vect[] 	Inventory vectors

type Inv struct {
	Type uint32 // 	Identifies the object type linked to this inventory
	Hash Hash   // 	Hash of the object
}

func MarshalInv(out []byte, v Inv) []byte {
	out = MarshalUint32(out, v.Type)
	out = MarshalHash(out, v.Hash)
	return out
}

func UnmarshalInv(data []byte) (Inv, []byte) {
	var v Inv
	v.Type, data = UnmarshalUint32(data)
	v.Hash, data = UnmarshalHash(data)
	return v, data
}

func MarshalInvs(out []byte, v []Inv) []byte {
	out = MarshalVarInt(out, uint64(len(v)))
	// can be made more efficient by preallocating space
	for _, h := range v {
		out = MarshalInv(out, h)
	}
	return out
}

func UnmarshalInvs(data []byte) ([]Inv, []byte) {
	l, data := UnmarshalVarInt(data)
	v := make([]Inv, l)
	for i := 0; i < int(l); i++ {
		v[i], data = UnmarshalInv(data)
	}
	return v, data
}

type InvMessage struct {
	Invs []Inv
}

func (msg InvMessage) Marshal(out []byte) []byte {
	return MarshalInvs(out, msg.Invs)
}

func (msg *InvMessage) Unmarshal(data []byte) []byte {
	msg.Invs, data = UnmarshalInvs(data)
	return data
}

func (msg InvMessage) GetCommandString() string {
	return "inv"
}
