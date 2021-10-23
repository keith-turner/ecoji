package ecoji

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	reader := strings.NewReader("🪐📩🎦🐇🛻📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🪳💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	buffer1 := &bytes.Buffer{}
	err := Decode(reader, buffer1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Base64 is so 1999, isn't there something better?\n"
	if string(buf) != expected {
		t.Fatalf("should decode to '%s', was: '%s'", expected, string(buf))
	}

}

func TestDecodeV1(t *testing.T) {
	reader := strings.NewReader("🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	buffer1 := &bytes.Buffer{}

	err := Decode(reader, buffer1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Base64 is so 1999, isn't there something better?\n"
	if string(buf) != expected {
		t.Fatalf("should decode to '%s', was: '%s'", expected, string(buf))
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
}
