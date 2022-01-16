package network

import (
	"reflect"
	"testing"
)

func TestDoubleHash(t *testing.T) {
	// from: https://en.bitcoin.it/wiki/Protocol_documentation#Common_standards
	// hello
	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824 (first round of sha-256)
	// 9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50 (second round of sha-256)
	digest := doubleHash([]byte("hello"))
	expect, _ := StringToHash("9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50")
	if digest != expect {
		t.Errorf("Incorrect double hashed value 'hello'->\n'%s' != \n'%s'", digest, expect)
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

func TestStringToHash(t *testing.T) {
	h, _ := StringToHash("AB0000")
	e := Hash([32]byte{0xAB, 0x00, 0x00})
	if !reflect.DeepEqual(h, e) {
		t.Errorf("Hashes don't match: %v!=%v", h, e)
	}
}
func TestRPCStringToHash(t *testing.T) {
	h, _ := RPCStringToHash("0000AB")
	e := Hash([32]byte{0xAB, 0x00, 0x00})
	if !reflect.DeepEqual(h, e) {
		t.Errorf("Hashes don't match: %v!=%v", h, e)
	}
}

func TestGetDifficulty(t *testing.T) {
	// Default difficulty
	d := GetDifficulty(0x1d00ffff)
	e := 1.0
	if !reflect.DeepEqual(d, e) {
		t.Errorf("Difficulties don't match: %v!=%v", d, e)
	}

	// Sample from here: https://en.bitcoin.it/wiki/Difficulty
	// Fixed missing digits at end (_2)
	d = GetDifficulty(0x1a44b9f2)
	e = 244112.4877743364_2
	if !reflect.DeepEqual(d, e) {
		t.Errorf("Difficulties don't match: %v!=%v", d, e)
	}

	// https://chainquery.com/bitcoin-cli/getblock
	// 00000000000000001e8d6829a8a21adc5d38d0a473b144b6765798e61f98bd1d
	d = GetDifficulty(0x1b0404cb)
	e = 16307.420938523983
	if !reflect.DeepEqual(d, e) {
		t.Errorf("Difficulties don't match: %v!=%v", d, e)
	}

}
