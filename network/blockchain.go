package network

import "time"

type Header struct {
	Version        uint32    // version 			int32_t 	Block version information (note, this is signed)
	PrevBlockHash  Hash      // prev_block 	char[32] 	The hash value of the previous block this particular block references
	MerkleRootHash Hash      // merkle_root 	char[32] 	The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp      time.Time // timestamp 		uint32_t 	A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits           Compact   // bits 				uint32_t 	The calculated difficulty target being used for this block
	Nonce          uint32    // nonce 				uint32_t 	The nonce used to generate this block… to allow variations of the header and compute different hashes
	// TxnCount       uint64    // txn_count 		var_int 	Number of transaction entries, this value is always 0
}

func HashHeader(h Header) Hash {
	data := []byte("hello")
	return doubleHash(data)
}

func MarshalHeader(out []byte, v Header) []byte {

	out = MarshalUint32(out, v.Version)
	out = MarshalHash(out, v.PrevBlockHash)
	out = MarshalHash(out, v.MerkleRootHash)
	out = MarshalTimestamp4(out, v.Timestamp)
	out = MarshalCompact(out, v.Bits)
	out = MarshalUint32(out, v.Nonce)
	return out
}

func UnmarshalHeader(data []byte) (Header, []byte) {
	var v Header
	v.Version, data = UnmarshalUint32(data)
	v.PrevBlockHash, data = UnmarshalHash(data)
	v.MerkleRootHash, data = UnmarshalHash(data)
	v.Timestamp, data = UnmarshalTimestamp4(data)
	v.Bits, data = UnmarshalCompact(data)
	v.Nonce, data = UnmarshalUint32(data)

	return v, data
}
