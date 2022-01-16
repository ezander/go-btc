package network

import (
	"reflect"
	"testing"
	"time"
)

func s2h(s string) Hash {
	digest, _ := StringToHash(s)
	return digest
}

func rs2h(s string) Hash {
	digest, _ := RPCStringToHash(s)
	return digest
}

func TestHashHeader1(t *testing.T) {
	// https://chainquery.com/bitcoin-cli/getblock
	// "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"
	// block height 1
	h := Header{
		Version:        1,
		PrevBlockHash:  rs2h("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
		MerkleRootHash: rs2h("0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098"),
		Timestamp:      time.Unix(1231469665, 0),
		Bits:           0x1d00ffff,
		Nonce:          2573394689,
	}
	out := MarshalHeader([]byte{}, h)
	hash := doubleHash(out)
	hashExpect := rs2h("00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048")
	if !reflect.DeepEqual(hash, hashExpect) {
		t.Errorf("Hashes did not match (\n%v (actual) != \n%v (expected))", hash, hashExpect)
	}
}

func TestHashHeader0(t *testing.T) {
	// https://chainquery.com/bitcoin-cli/getblock
	// "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	// block height 0 (genesis block)
	h := Header{
		Version:        1,
		PrevBlockHash:  Hash{},
		MerkleRootHash: rs2h("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
		Timestamp:      time.Unix(1231006505, 0),
		Bits:           0x1d00ffff,
		Nonce:          2083236893,
	}
	out := MarshalHeader([]byte{}, h)
	hash := doubleHash(out)
	hashExpect := rs2h("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	if !reflect.DeepEqual(hash, hashExpect) {
		t.Errorf("Hashes did not match (\n%v (actual) != \n%v (expected))", hash, hashExpect)
	}
}
