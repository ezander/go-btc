package network

import (
	"reflect"
	"testing"
)

func TestReversed(t *testing.T) {
	s := []byte{1, 2, 3, 5}
	r := reversed(s)
	e := []byte{5, 3, 2, 1}
	if !reflect.DeepEqual(r, e) {
		t.Errorf("Not reversed '%v'->'%v' ", s, r)
	}
	if s[0] != 1 {
		t.Errorf("Input should not change '%v'", s)
	}
}

func TestReverse(t *testing.T) {
	s := []byte{1, 2, 3, 5}
	r := reverse(s)
	e := []byte{5, 3, 2, 1}
	if !reflect.DeepEqual(r, e) {
		t.Errorf("Not reversed '%v'->'%v' ", s, r)
	}
	if !reflect.DeepEqual(s, e) {
		t.Errorf("Input should change '%v'", s)
	}
}

func TestReversedAsciiString(t *testing.T) {
	s := "hello World!"
	r := reversedAsciiString(s)
	e := "!dlroW olleh"
	if !reflect.DeepEqual(r, e) {
		t.Errorf("Not reversed '%v'->'%v' ", s, r)
	}
	if s[0] != 'h' {
		t.Errorf("Input should not change '%v'", s)
	}
}

func TestReversedHexString(t *testing.T) {
	s := "12ab34f6"
	r := reversedHexString(s)
	e := "f634ab12"
	if !reflect.DeepEqual(r, e) {
		t.Errorf("Not reversed '%v'->'%v' ", s, r)
	}
	if s[0] != '1' {
		t.Errorf("Input should not change '%v'", s)
	}
}
