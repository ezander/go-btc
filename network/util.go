package network

import "encoding/hex"

// Reverse reverses the order in a byte array (modifies the input)
func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func reversed(s []byte) []byte {
	r := make([]byte, len(s))
	for i, j := 0, len(s)-1; i < len(s); i, j = i+1, j-1 {
		r[i] = s[j]
	}
	return r
}

func reversedAsciiString(s string) string {
	b := []byte(s)
	return string(reverse(b))
}

func reversedHexString(s string) string {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(reverse(b))
}
