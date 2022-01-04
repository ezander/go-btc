package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

// Interface for marshalling
type Marshaller interface {
	Marshal(out []byte) []byte
}

type Unmarshaller interface {
	Unmarshal(data []byte) []byte
}

// Helper functions for marshalling and unmarshalling integers
func MarshalUint8(out []byte, v uint8) []byte {
	return append(out, v)
}

func UnmarshalUint8(data []byte) (uint8, []byte) {
	return data[0], data[1:]
}

func MarshalUint16(out []byte, v uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, v)
	return append(out, data...)
}

func UnmarshalUint16(data []byte) (uint16, []byte) {
	return binary.LittleEndian.Uint16(data), data[2:]
}

func MarshalUint32(out []byte, v uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, v)
	return append(out, data...)
}

func UnmarshalUint32(data []byte) (uint32, []byte) {
	return binary.LittleEndian.Uint32(data), data[4:]
}

func MarshalUint64(out []byte, v uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, v)
	return append(out, data...)
}

func UnmarshalUint64(data []byte) (uint64, []byte) {
	return binary.LittleEndian.Uint64(data), data[8:]
}

// Other integer types
func MarshalVarInt(out []byte, v uint64) []byte {
	switch {
	case v < 0xFD:
		return MarshalUint8(out, uint8(v))
	case v <= 0xFFFF:
		return MarshalUint16(MarshalUint8(out, 0xFD), uint16(v))
	case v <= 0xFFFFFFFF:
		return MarshalUint32(MarshalUint8(out, 0xFE), uint32(v))
	default:
		return MarshalUint64(MarshalUint8(out, 0xFF), uint64(v))
	}
}

func UnmarshalVarInt(data []byte) (uint64, []byte) {
	b := data[0]
	if b < 0xFD {
		value, data := UnmarshalUint8(data)
		return uint64(value), data
	} else {
		switch b {
		case 0xFD:
			value, data := UnmarshalUint16(data[1:])
			return uint64(value), data
		case 0xFE:
			value, data := UnmarshalUint32(data[1:])
			return uint64(value), data
		default:
			value, data := UnmarshalUint64(data[1:])
			return uint64(value), data
		}
	}
}

func MarshalBool(out []byte, v bool) []byte {
	if v {
		return MarshalUint8(out, 1)
	} else {
		return MarshalUint8(out, 0)
	}
}

func UnmarshalBool(data []byte) (bool, []byte) {
	b, data := UnmarshalUint8(data)
	return b != 0, data
}

// Timestamps
func MarshalTimestamp(out []byte, v time.Time) []byte {
	return MarshalUint64(out, uint64(v.Unix()))
}

func UnmarshalTimestamp(data []byte) (time.Time, []byte) {
	value, data := UnmarshalUint64(data)
	return time.Unix(int64(value), 0), data
}

// todo: Marshal and Unmarshal Timestamp4

// String types

func MarshalVarStr(out []byte, v string) []byte {
	out = MarshalVarInt(out, uint64(len(v)))
	out = append(out, []byte(v)...)
	return out
}

func UnmarshalVarStr(data []byte) (string, []byte) {
	l, data := UnmarshalVarInt(data)
	return string(data[:l]), data[l:]
}

func MarshalFixedStr(out []byte, v string, l int) []byte {
	s := make([]byte, l)
	copy(s[:], []byte(v))
	return append(out, s[:]...)
}

func UnmarshalFixedStr(data []byte, l int) (string, []byte) {
	v := string(data[:l])
	v = strings.TrimRight(v, "\x00")
	return v, data[l:]
}

func MarshalBytes(out []byte, v []byte) []byte {
	return append(out, v...)
}

func UnmarshalBytes(data []byte, l uint32) ([]byte, []byte) {
	return data[:l], data[l:]
}

// Network related types
func MarshalIP(out []byte, v net.IP) []byte {
	bytes := []byte(v)
	if len(bytes) != 16 {
		fmt.Println(bytes)
		fmt.Println(v.String())
		panic(fmt.Sprintf("Wrong length of IPv6 address (was %d, expected 16)", len(bytes)))
	}
	return MarshalBytes(out, bytes)
}

func UnmarshalIP(out []byte) (net.IP, []byte) {
	bytes, out := UnmarshalBytes(out, 16)
	return net.IP(bytes), out
}

type NetAddr struct {
	Services uint64
	IPAddr   net.IP
	Port     uint16
}

func MarshalNetAddr(out []byte, v NetAddr) []byte {
	out = MarshalUint64(out, v.Services)
	out = MarshalIP(out, v.IPAddr)
	out = MarshalUint16(out, v.Port)
	return out
}

func UnmarshalNetAddr(out []byte) (NetAddr, []byte) {
	var v NetAddr
	v.Services, out = UnmarshalUint64(out)
	v.IPAddr, out = UnmarshalIP(out)
	v.Port, out = UnmarshalUint16(out)
	return v, out
}

// type TimestampedNetAddr struct {
// 	Time time.Time
// 	NetAddr
// }

// func (v TimestampedNetAddr) Marshal(out []byte) []byte {
// 	out = MarshalTimestamp(out, v.Time)
// 	out = v.NetAddr.Marshal(out)
// 	return out
// }
