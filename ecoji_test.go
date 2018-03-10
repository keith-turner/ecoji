package ecoji

import (
	"bytes"
	"strings"
	"testing"
)

func TestOneByteEncode(t *testing.T) {
	reader := strings.NewReader("t")
	buffer1 := bytes.NewBuffer(nil)

	Encode(reader, buffer1, 0)
	actual, _ := buffer1.ReadString('\n')

	expected := []rune{mapping[int('t')<<2], padding, padding, padding}

	if cmp := strings.Compare(actual, string(expected)); cmp != 0 {
		t.Errorf("%s != %s %d", string(expected), actual, cmp)
	}
}

func TestGarbage(t *testing.T) {
	reader := strings.NewReader("not emojis")
	buffer1 := bytes.NewBuffer(nil)

	err := Decode(reader, buffer1)
	if err == nil {
		t.Error("Expected error")
	}

	if !strings.Contains(err.Error(), "Invalid rune") {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}
