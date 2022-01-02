package network

import (
	"crypto/sha256"
	"encoding/binary"
)

func doubleHash(data []byte) [32]byte {
	digest1 := sha256.Sum256(data)
	digest2 := sha256.Sum256(digest1[:])
	return digest2
}

func checksum(data []byte) uint32 {
	digest := doubleHash(data)
	return binary.LittleEndian.Uint32(digest[:4])
}
