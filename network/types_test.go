package network

import (
	"testing"
	"time"
)

func TestMarshalUint64(t *testing.T) {
	var x uint64 = 1234
	data := marshalUint64([]byte{}, x)
	y, data := unmarshalUint64(data)
	if len(data) > 0 {
		t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
	}
	if y != x {
		t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
	}
}

func TestMarshalVarInt(t *testing.T) {
	vals := []uint64{1, 20, 0xFF, 0xFAAF, 0xFAAFFBBF, 0x1FAAFFBBF}
	lens := []int{1, 1, 3, 3, 5, 9}
	for i, x := range vals {
		data := MarshalVarInt([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for 0x%X (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalVarInt(data)
		if len(data) > 0 {
			t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalTimestamp(t *testing.T) {
	x := time.Now().Truncate(time.Second)
	data := MarshalTimestamp([]byte{}, x)
	y, data := UnmarshalTimestamp(data)
	if len(data) > 0 {
		t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
	}
	if y != x { // Exact equality only because of truncating to seconds in initialization
		t.Errorf("Unmarshalled data did not match (%v!=%v)", time.Time(y).String(), x)
	}
}

func TestMarshalVarString(t *testing.T) {
	vals := []string{"", "foo", "lakdslfjsldf"}
	for _, x := range vals {
		data := MarshalVarStr([]byte{}, x)
		y, data := UnmarshalVarStr(data)
		if len(data) > 0 {
			t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalFixedString(t *testing.T) {
	vals := []string{"", "foo", "012345678901"}
	for _, x := range vals {
		data := MarshalFixedStr([]byte{}, x, 12)
		if len(data) != 12 {
			t.Errorf("Incorrect length of marshalled data for string %s (%v!=12)", x, len(data))
		}
		y, data := UnmarshalFixedStr(data, 12)
		if len(data) > 0 {
			t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}
