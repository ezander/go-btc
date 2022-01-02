package network

import (
	"fmt"
	"testing"
)

func TestDoubleHash(t *testing.T) {
	// from: https://en.bitcoin.it/wiki/Protocol_documentation#Common_standards
	// hello
	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824 (first round of sha-256)
	// 9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50 (second round of sha-256)
	digest := doubleHash([]byte("hello"))
	s := fmt.Sprintf("%x", digest)
	if s != "9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50" {
		t.Errorf("Incorrect double hashed value 'hello'->'%s' ", s)
	}
}

func TestChecksum(t *testing.T) {
	// from: https://en.bitcoin.it/wiki/Protocol_documentation#Common_standards
	// hello
	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824 (first round of sha-256)
	// 9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50 (second round of sha-256)
	// take first 4 bytes (8 nibbles) in reversed order (i.e. little endian)
	cs := checksum([]byte("hello"))
	if cs != 0xdfc99595 {
		t.Errorf("Incorrect checksum value 'hello'->'%x' ", cs)
	}
}
