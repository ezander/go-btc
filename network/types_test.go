package network

import (
	"net"
	"testing"
	"time"
)

func TestMarshalUint8(t *testing.T) {
	vals := []uint8{0x00, 0xFE}
	lens := []int{1, 1}
	for i, x := range vals {
		data := MarshalUint8([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for 0x%X (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalUint8(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalUint16(t *testing.T) {
	vals := []uint16{0x0000, 0xFAFB}
	lens := []int{2, 2}
	for i, x := range vals {
		data := MarshalUint16([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for 0x%X (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalUint16(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalUint32(t *testing.T) {
	vals := []uint32{0x00000000, 0xFFFEFDFC}
	lens := []int{4, 4}
	for i, x := range vals {
		data := MarshalUint32([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for 0x%X (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalUint32(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalUint64(t *testing.T) {
	vals := []uint64{0x0000000000000000, 0xF0F1F2F3F4F5F6F7}
	lens := []int{8, 8}
	for i, x := range vals {
		data := MarshalUint64([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for 0x%X (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalUint64(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalVarInt(t *testing.T) {
	vals := []uint64{1, 20, 0xFF, 0xFFAF, 0xFFAFFBBF, 0xFFAAFFBBF}
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

func TestMarshalBool(t *testing.T) {
	vals := []bool{true, false}
	lens := []int{1, 1}
	for i, x := range vals {
		data := MarshalBool([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalBool(data)
		if len(data) > 0 {
			t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalTimestamp(t *testing.T) {
	vals := []time.Time{time.Now(), time.Unix(0, 0), time.Unix(1_000_000_000, 0)}
	lens := []int{8, 8, 8}
	for i, x := range vals {
		// Need to truncate to second, because fractions of seconds aren't transmitted in the BTC protocol
		x = x.Truncate(time.Second)
		data := MarshalTimestamp([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalTimestamp(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y.String(), x.String())
		}
	}
}

func TestMarshalTimestamp4(t *testing.T) {
	vals := []time.Time{time.Now(), time.Unix(0, 0), time.Unix(1_000_000_000, 0)}
	lens := []int{4, 4, 4}
	for i, x := range vals {
		// Need to truncate to second, because fractions of seconds aren't transmitted in the BTC protocol
		x = x.Truncate(time.Second)
		data := MarshalTimestamp4([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalTimestamp4(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y != x {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y.String(), x.String())
		}
	}
}

func TestMarshalVarString(t *testing.T) {
	vals := []string{"", "foo", "abcdefghij", string(make([]byte, 0xFF)), string(make([]byte, 0x10000))}
	lens := []int{1, 3 + 1, 10 + 1, 0xFF + 3, 0x10000 + 5}
	for i, x := range vals {
		data := MarshalVarStr([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", "", len(data), lens[i])
		}
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

func TestMarshalBytes(t *testing.T) {
	vals := [][]byte{{}, {1, 2, 3}}
	lens := []int{0, 3}
	for i, x := range vals {
		data := MarshalBytes([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", "", len(data), lens[i])
		}
		y, data := UnmarshalBytes(data, uint32(lens[i]))
		if len(data) > 0 {
			t.Errorf("Len Data should be zero after unmarshalling (%d)", len(data))
		}
		if string(y) != string(x) {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y, x)
		}
	}
}

func TestMarshalIP(t *testing.T) {
	vals := []net.IP{net.IPv4zero, net.IPv4(127, 0, 0, 1), net.IPv6zero, net.IPv6loopback}
	lens := []int{16, 16, 16, 16}
	for i, x := range vals {
		data := MarshalIP([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalIP(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if string(y) != string(x) {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y.String(), x.String())
		}
	}
}

func TestNetAddr(t *testing.T) {
	vals := []NetAddr{{1234, net.IPv4zero, 8333}, {34, net.IPv6loopback, 18333}}
	lens := []int{26, 26} // 8 + 16 + 2
	for i, x := range vals {
		data := MarshalNetAddr([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalNetAddr(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y.String() != x.String() {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y.String(), x.String())
		}
	}
}

func TestTimeNetAddr(t *testing.T) {
	vals := []TimeNetAddr{{time.Unix(0, 0), NetAddr{1234, net.IPv4zero, 8333}}, {time.Unix(1_000_000_000, 0), NetAddr{34, net.IPv6loopback, 18333}}}
	lens := []int{30, 30} // 4 + 8 + 16 + 2
	for i, x := range vals {
		data := MarshalTimeNetAddr([]byte{}, x)
		if len(data) != lens[i] {
			t.Errorf("Incorrect length of marshalled data for '%v' (%v!=%v)", x, len(data), lens[i])
		}
		y, data := UnmarshalTimeNetAddr(data)
		if len(data) > 0 {
			t.Errorf("Len of data should be zero after unmarshalling (%d)", len(data))
		}
		if y.String() != x.String() {
			t.Errorf("Unmarshalled data did not match (%v!=%v)", y.String(), x.String())
		}
	}
}
