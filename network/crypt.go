package network

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Hash [32]byte

func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

func StringToHash(s string) (Hash, error) {
	var h Hash
	b, error := hex.DecodeString(s)
	if error != nil {
		return h, error
	}
	copy(h[:], b[:])
	return h, nil
}

func doubleHash(data []byte) Hash {
	digest1 := sha256.Sum256(data)
	digest2 := sha256.Sum256(digest1[:])
	return digest2
}

func checksum(data []byte) uint32 {
	digest := doubleHash(data)
	return binary.LittleEndian.Uint32(digest[:4])
}
