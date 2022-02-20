package ecoji

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func decode(s string) (string, error) {
	reader := strings.NewReader(s)
	buffer1 := &bytes.Buffer{}
	err := Decode(reader, buffer1)
	if err != nil {
		return "", err
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func TestDecodeConcatenated(t *testing.T) {
	dstr, err := decode("👖📸🧊🌭👩☕💲🥇🪚☕")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "abcdefxyz"
	if dstr != expected {
		t.Fatalf("should decode to '%s', was: '%s'", expected, dstr)
	}
}

func TestDecodeMixed(t *testing.T) {

	// the 2nd rune is ecoji v1 only and the 3rd rune is ecoji v2 only
	runes := [4]rune{0x1f004, 0x1f170, 0x1f93f, 0x1f93e}

	reader := strings.NewReader(string(runes[:]))
	buffer1 := &bytes.Buffer{}

	err := Decode(reader, buffer1)

	if err == nil {
		t.Errorf("Did not see error with mixed data")
	} else if !strings.Contains(err.Error(), "Emojis from different ecoji versions seen") {
		t.Errorf("Did not see expected error message")
	}

	// the 2nd rune is ecoji v2 only and the 3rd rune is ecoji v1 only
	runes2 := [4]rune{0x1f004, 0x1f93f, 0x1f170, 0x1f93e}

	reader2 := strings.NewReader(string(runes2[:]))
	buffer2 := &bytes.Buffer{}

	err2 := Decode(reader2, buffer2)

	if err2 == nil {
		t.Errorf("Did not see error with mixed data")
	} else if !strings.Contains(err2.Error(), "Emojis from different ecoji versions seen") {
		t.Errorf("Did not see expected error message")
	}

	// test where mixed are more than 4 apart
	runes3 := []rune{0x1f004, 0x1f170, 0x1f170, 0x1f93e, 0x1f004, 0x1f170, 0x1f170, 0x1f93e, 0x1f004, 0x1f170, 0x1f93f, 0x1f93e}

	reader3 := strings.NewReader(string(runes3))
	buffer3 := &bytes.Buffer{}

	err3 := Decode(reader3, buffer3)

	if err3 == nil {
		t.Errorf("Did not see error with mixed data")
	} else if !strings.Contains(err3.Error(), "Emojis from different ecoji versions seen") {
		t.Errorf("Did not see expected error message")
	}

}
