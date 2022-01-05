package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"
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

func PublicKeyFromString(s string) ecdsa.PublicKey {
	reverse := func(s []byte) []byte {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		return s
	}

	buffer, err := hex.DecodeString(s)
	if err != nil {
		panic("Could not decode public key")
	}
	if buffer[0] != 0x04 {
		panic("Public key representation should start with 0x04")
	}
	if len(buffer) != 1+2*32 {
		panic("Public key reprentation has wrong length")
	}
	toBigInt := func(n []byte) *big.Int {
		X := big.NewInt(0)
		X.SetBytes(n)
		if !true {
			X.SetBytes(reverse(n))
		}
		return X
	}
	X := toBigInt(buffer[1:33])
	Y := toBigInt(buffer[33:65])
	curve := elliptic.P256()
	pub := ecdsa.PublicKey{Curve: curve, X: X, Y: Y}
	fmt.Println("pub: ", pub)
	return pub
}

func VerifyASN1(pub *ecdsa.PublicKey, hash, sig []byte) bool {
	// shamelessly copied from go1.15
	var (
		r, s  = &big.Int{}, &big.Int{}
		inner cryptobyte.String
	)
	input := cryptobyte.String(sig)
	if !input.ReadASN1(&inner, asn1.SEQUENCE) ||
		!input.Empty() ||
		!inner.ReadASN1Integer(r) ||
		!inner.ReadASN1Integer(s) ||
		!inner.Empty() {
		return false
	}
	return ecdsa.Verify(pub, hash, r, s)
}
