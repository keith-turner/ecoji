package ecoji

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	reader := strings.NewReader("🦥📩🎦🐇🧖📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🥲💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	buffer1 := &bytes.Buffer{}

	err := Decode(reader, buffer1)
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	expected := "Base64 is so 1999, isn't there something better?\n"
	if string(buf) != expected {
		t.Fatalf("should decode to '%s', was: %s", expected, string(buf))
	}

}

func TestDecodeV1(t *testing.T) {
	reader := strings.NewReader("🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	buffer1 := &bytes.Buffer{}

	err := Decode(reader, buffer1)
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if string(buf) != "Base64 is so 1999, isn't there something better?\n" {
		t.Fatal("should decode to abc was:", string(buf))
	}
}
