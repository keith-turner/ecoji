package ecoji

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func writeBytes(t *testing.T, data []byte, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for _, b := range data {
		writer.WriteByte(b)
	}
}

func writeRunes(t *testing.T, data []rune, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for _, r := range data {
		writer.WriteRune(r)
	}
}

func check(t *testing.T, expectedV1 []rune, expectedV2 []rune, input []byte, name string, wrap uint) {
	reader := bytes.NewBuffer(input)
	buffer1 := bytes.NewBuffer(nil)

	err := EncodeV2(reader, buffer1, wrap)
	if err != nil {
		t.Error(err)
	}

	actual := buffer1.String()

	if cmp := strings.Compare(actual, string(expectedV2)); cmp != 0 {
		t.Errorf("'%s' != '%s' %d", string(expectedV2), actual, cmp)
	}

	buffer2 := bytes.NewBuffer(nil)
	err = Decode(strings.NewReader(string(expectedV2)), buffer2)
	if err != nil {
		t.Error(err)
	}

	if cmp := bytes.Compare(input, buffer2.Bytes()); cmp != 0 {
		t.Errorf("'%v' != '%v' %d", input, buffer2.Bytes(), cmp)
	}

	// start checking V1
	reader2 := bytes.NewBuffer(input)
	buffer3 := bytes.NewBuffer(nil)

	err = Encode(reader2, buffer3, wrap)
	if err != nil {
		t.Error(err)
	}

	actual = buffer3.String()

	if cmp := strings.Compare(actual, string(expectedV1)); cmp != 0 {
		t.Errorf("'%s' != '%s' %d", string(expectedV1), actual, cmp)
	}

	buffer4 := bytes.NewBuffer(nil)
	err = Decode(strings.NewReader(string(expectedV1)), buffer4)
	if err != nil {
		t.Error(err)
	}

	if cmp := bytes.Compare(input, buffer4.Bytes()); cmp != 0 {
		t.Errorf("'%v' != '%v' %d", input, buffer4.Bytes(), cmp)
	}

	if name != "" {
		writeBytes(t, input, "../test_scripts/data/"+name+".plain")
		writeRunes(t, expectedV1, "../test_scripts/data/"+name+".ev1")
		writeRunes(t, expectedV2, "../test_scripts/data/"+name+".ev2")
	}
}

func TestZeroByteEncode(t *testing.T) {
	check(t, []rune{}, []rune{}, []byte{}, "zero_byte", 0)
}

func TestOneByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[int('k')<<2], padding, padding, padding}, []rune{emojisV2[int('k')<<2], padding}, []byte{'k'}, "one_byte", 0)
}

func TestTwoByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], padding, padding}, []rune{emojisV2[0], emojisV2[16], padding}, []byte{0x00, 0x01}, "two_byte", 0)
}

func TestThreeByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], padding}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], padding}, []byte{0x00, 0x01, 0x02}, "three_byte", 0)
}

func TestFourByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[0]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[0]}, []byte{0x00, 0x01, 0x02, 0x00}, "four_byte_1", 0)
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[1]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[1]}, []byte{0x00, 0x01, 0x02, 0x01}, "four_byte_2", 0)
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[2]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[2]}, []byte{0x00, 0x01, 0x02, 0x02}, "four_byte_3", 0)
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[3]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[3]}, []byte{0x00, 0x01, 0x02, 0x03}, "four_byte_4", 0)
}

func TestFiveByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[687], emojisV1[222], emojisV1[960], emojisV1[291]}, []rune{emojisV2[687], emojisV2[222], emojisV2[960], emojisV2[291]}, []byte{0xab, 0xcd, 0xef, 0x01, 0x23}, "five_byte", 0)
}

func TestSixByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 123<<38 | 456<<28 | 789<<18 | 909<<8 | 55
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[123], emojisV1[456], emojisV1[789], emojisV1[909], emojisV1[55<<2], padding, padding, padding},
		[]rune{emojisV2[123], emojisV2[456], emojisV2[789], emojisV2[909], emojisV2[55<<2], padding}, data[2:8], "six_byte", 0)
}

func TestSevenByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 237<<46 | 77<<36 | 257<<26 | 513<<16 | 809<<6 | 7
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[237], emojisV1[77], emojisV1[257], emojisV1[513], emojisV1[809], emojisV1[7<<4], padding, padding},
		[]rune{emojisV2[237], emojisV2[77], emojisV2[257], emojisV2[513], emojisV2[809], emojisV2[7<<4], padding}, data[1:8], "seven_byte", 0)
}

func TestEightByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 3<<54 | 206<<44 | 368<<34 | 617<<24 | 650<<14 | 1005<<4 | 3
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[3], emojisV1[206], emojisV1[368], emojisV1[617], emojisV1[650], emojisV1[1005], emojisV1[3<<6], padding},
		[]rune{emojisV2[3], emojisV2[206], emojisV2[368], emojisV2[617], emojisV2[650], emojisV2[1005], emojisV1[3<<6], padding}, data[0:8], "eight_byte", 0)
}

func TestNineByteEncode(t *testing.T) {

	for i := 0; i < 4; i++ {
		var scratch uint64
		scratch = (855<<30 | 298<<20 | 1007<<10 | 97) << 24
		var data [13]byte
		binary.BigEndian.PutUint64(data[0:8], scratch)

		scratch = (611<<30 | 291<<20 | 856<<10 | uint64(i)<<8) << 24
		binary.BigEndian.PutUint64(data[5:13], scratch)

		check(t, []rune{emojisV1[855], emojisV1[298], emojisV1[1007], emojisV1[97], emojisV1[611], emojisV1[291], emojisV1[856], paddingLastV1[i]},
			[]rune{emojisV2[855], emojisV2[298], emojisV2[1007], emojisV2[97], emojisV2[611], emojisV2[291], emojisV2[856], paddingLastV2[i]}, data[0:9], "nine_byte_"+fmt.Sprint(i), 0)
	}
}

func TestExhaustive(t *testing.T) {
	// use this to hold 10 bit ordinals
	biggy := big.NewInt(1)

	myRand := rand.New(rand.NewSource(42))

	// create an array w/ all 1024 ordinals in random order
	ordinals := myRand.Perm(1024)

	var expectedRunesV1 [1024]rune
	var expectedRunesV2 [1024]rune

	for i, _ := range ordinals {
		expectedRunesV1[i] = emojisV1[ordinals[i]]
		expectedRunesV2[i] = emojisV2[ordinals[i]]

		// shift left by to bits and then add an ordinal
		biggy.Lsh(biggy, 10)
		biggy.Add(biggy, big.NewInt((int64)(ordinals[i])))
	}

	// get the ordinals encoded as 10 bit integers, ignoring the 1st byte as it contains the initial 1 used to create
	// biggy
	bytes := biggy.Bytes()[1:]

	check(t, expectedRunesV1[:], expectedRunesV2[:], bytes, "exhaustive", 0)
}

func TestPhrase(t *testing.T) {
	expectedV1 := []rune("🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	expectedV2 := []rune("🧏📩🧈🐇🧅📘🔯🚜💞😽♏🐊🎱🥁🚄🌱💞😭💮✊💢🪠🐭🩴🍉🚲🦑🐶💢🪠🔮🩹🍉📸🐮🌼👦🚟🥴📑")
	plain := []byte("Base64 is so 1999, isn't there something better?\n")
	check(t, expectedV1, expectedV2, plain, "phrase", 0)
}

func TestWrap(t *testing.T) {
	check(t, []rune("🎌🚟🎗🈸🎥🤠📠🐁👖📸🎈☕"), []rune("🎌🚟🦿🦣🎥🤠📠🐁👖📸🎈☕"), []byte("1234567890abc"), "", 0)
	check(t, []rune("🎌\n🚟\n🎗\n🈸\n🎥\n🤠\n📠\n🐁\n👖\n📸\n🎈\n☕\n"),
		[]rune("🎌\n🚟\n🦿\n🦣\n🎥\n🤠\n📠\n🐁\n👖\n📸\n🎈\n☕\n"),
		[]byte("1234567890abc"), "", 1)
	check(t, []rune("🎌🚟\n🎗🈸\n🎥🤠\n📠🐁\n👖📸\n🎈☕\n"),
		[]rune("🎌🚟\n🦿🦣\n🎥🤠\n📠🐁\n👖📸\n🎈☕\n"),
		[]byte("1234567890abc"), "", 2)
	check(t, []rune("🎌🚟🎗\n🈸🎥🤠\n📠🐁👖\n📸🎈☕\n"),
		[]rune("🎌🚟🦿\n🦣🎥🤠\n📠🐁👖\n📸🎈☕\n"),
		[]byte("1234567890abc"), "", 3)
	check(t, []rune("🎌🚟🎗🈸\n🎥🤠📠🐁\n👖📸🎈☕\n"),
		[]rune("🎌🚟🦿🦣\n🎥🤠📠🐁\n👖📸🎈☕\n"),
		[]byte("1234567890abc"), "", 4)
	check(t, []rune("🎌🚟🎗🈸🎥\n🤠📠🐁👖📸\n🎈☕\n"),
		[]rune("🎌🚟🦿🦣🎥\n🤠📠🐁👖📸\n🎈☕\n"),
		[]byte("1234567890abc"), "", 5)
	check(t, []rune("🎌🚟🎗🈸🎥🤠📠🐁👖📸🎈☕\n"),
		[]rune("🎌🚟🦿🦣🎥🤠📠🐁👖📸🎈☕\n"),
		[]byte("1234567890abc"), "", 20)

}

func TestWindowsNewLine(t *testing.T) {
	testDecode(t, "🎌\r\n🚟\r\n🦿🦣🎥🤠\r\n📠🐁👖📸🎈☕", []byte("1234567890abc"), "windows_newline_v2_1")
	testDecode(t, "🎌🚟🎗🈸🎥🤠📠🐁👖📸🎈☕\r\n", []byte("1234567890abc"), "windows_newline_v2_2")
}

func decode(s string) ([]byte, error) {
	reader := strings.NewReader(s)
	buffer1 := &bytes.Buffer{}
	err := Decode(reader, buffer1)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(buffer1)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func testDecode(t *testing.T, encoded string, expected []byte, name string) {
	dstr, err := decode(encoded)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cmp := bytes.Compare(expected, dstr); cmp != 0 {
		t.Fatalf("should decode to '%s', was: '%s'", expected, dstr)
	}

	if name != "" {
		writeBytes(t, expected, "../test_scripts/data/"+name+".plaind")
		writeRunes(t, []rune(encoded), "../test_scripts/data/"+name+".enc")
	}
}

func TestDecodeConcatenated(t *testing.T) {
	// V2 Encoded and concat
	testDecode(t, "👖📸🧈🌭👩☕💲🥇🪚☕", []byte("abcdefxyz"), "concat_v2_1")
	// V1 Encoded and concat
	testDecode(t, "👖📸🎈☕🎥🤠📠🏍🐲👡🕟☕", []byte("abc6789XY\n"), "concat_v1_1")

	// Test V1 concat of encoded messages of lengths 1 to 10.  So did enc("A")+enc("BC")+enc("DEF")+...+enc("jklmnopqrs")
	testDecode(t, "🏒☕☕☕🏗🈳☕☕🏟🌚👑☕🏫🍌🔥📑🏾🎌🛡🔢🐒🏣🍜🛢🐥☕☕☕🐪👆📨🐫🎈🚌☕☕🎐🚯🏛🐇🎩🤰🔓☕👖📸🎦🌭👪🕕📬🏍👺😁🚗🐿💎🚃🌤🕒",
		[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrs"), "concat_v1_2")

	// Test V2 concat of encoded messages of lengths 1 to 10.
	testDecode(t, "🏒☕🧏🥱☕🧝🌚👑☕🏫🍌🔥📑🧦🎌🫣🧽🐒🏣🍜🫤🐥☕🐪👆📨🐫🎈🚌☕🎐🚯🧙🐇🎩🤰🔓☕👖📸🧈🌭👪🪐📬🛼👺😁🚗🧨💎🚃🦩🪄",
		[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrs"), "concat_v2_2")

}

func testGarbageInput(t *testing.T, testData string, expectedErrMsg string, outputFileName string) {
	reader := strings.NewReader(testData)
	buffer1 := bytes.NewBuffer(nil)

	err := Decode(reader, buffer1)
	if err == nil {
		t.Error("Expected error for : " + outputFileName)
		return
	}

	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Unexpected error message: '%s'  expected '%s'", err.Error(), expectedErrMsg)
		return
	}

	writeRunes(t, []rune(testData), "../test_scripts/data/"+outputFileName+".garbage")

}

func TestGarbage(t *testing.T) {
	testGarbageInput(t, "not emojisV2", "Non Ecoji character seen", "ascii")
	// test emojis not used by ecoji
	testGarbageInput(t, "🟠🟡🤍🟩", "Non Ecoji character seen", "non_ecoji_emoji")
	// test ecoji v1 data that does not have expected padding.  The emojis below are only used in ecoji v1
	testGarbageInput(t, "🌶🌶🌶🌶🌶", "Unexpected end of data, input data size not multiple of 4", "incorrect_v1_padding")
	// test misplaced padding
	runes1 := [4]rune{paddingLastV1[0], emojisV1[1], emojisV1[2], emojisV1[3]}
	testGarbageInput(t, string(runes1[:]), "Last padding seen in unexpected position", "misplaced_padding_1")
	runes2 := [4]rune{paddingLastV2[0], emojisV2[1], emojisV2[2], emojisV2[3]}
	testGarbageInput(t, string(runes2[:]), "Last padding seen in unexpected position", "misplaced_padding_2")
	runes3 := [4]rune{emojisV1[1], paddingLastV1[1], emojisV1[2], emojisV1[3]}
	testGarbageInput(t, string(runes3[:]), "Last padding seen in unexpected position", "misplaced_padding_3")
	runes4 := [4]rune{emojisV2[1], emojisV2[2], paddingLastV2[3], emojisV2[3]}
	testGarbageInput(t, string(runes4[:]), "Last padding seen in unexpected position", "misplaced_padding_4")

	//test an ecoji v1 data that must padd out to len 4 and does not
	runes5 := [4]rune{0x1f004, 0x1f170, padding, 0x1f93e}
	testGarbageInput(t, string(runes5[:]), "Unexpectedly saw non-padding after padding", "misplaced_padding_5")

	//fill padding cannot be first char in 4 char block
	runes6 := [4]rune{padding, emojisV2[1], emojisV2[2], emojisV2[3]}
	testGarbageInput(t, string(runes6[:]), "Padding unexpectedly seen in first position", "misplaced_padding_6")
	runes7 := [8]rune{emojisV2[1], emojisV2[2], emojisV2[3], emojisV2[4], padding, padding, padding, padding}
	testGarbageInput(t, string(runes7[:]), "Padding unexpectedly seen in first position", "misplaced_padding_7")

	//test input data that is not a multiple of 4 and does not end with padding
	runes8 := [3]rune{emojisV2[1], emojisV2[2], emojisV2[3]}
	testGarbageInput(t, string(runes8[:]), "Unexpected end of data, input data size not multiple of 4", "missing_padding_1")
	runes9 := [5]rune{emojisV2[1], emojisV2[2], emojisV2[3], emojisV2[4], emojisV2[5]}
	testGarbageInput(t, string(runes9[:]), "Unexpected end of data, input data size not multiple of 4", "missing_padding_2")

	testGarbageInput(t, "🎌\r🚟\r🎗🈸🎥🤠\r📠🐁👖📸🎈☕", "Saw \r that was not followed by \n", "bad_newline_1")
	testGarbageInput(t, "🎌🚟🎗🈸🎥🤠📠🐁👖📸🎈☕\r", "Saw \r that was not followed by \n", "bad_newline_2")

}

func TestDecodeMixed(t *testing.T) {

	// the 2nd rune is ecoji v1 only and the 3rd rune is ecoji v2 only
	runes := [4]rune{0x1f004, 0x1f170, 0x1f93f, 0x1f93e}
	testGarbageInput(t, string(runes[:]), "Emojis from different ecoji versions seen", "mixed_1")

	// the 2nd rune is ecoji v2 only and the 3rd rune is ecoji v1 only
	runes2 := [4]rune{0x1f004, 0x1f93f, 0x1f170, 0x1f93e}
	testGarbageInput(t, string(runes2[:]), "Emojis from different ecoji versions seen", "mixed_2")

	// test where mixed are more than 4 apart
	runes3 := []rune{0x1f004, 0x1f170, 0x1f170, 0x1f93e, 0x1f004, 0x1f170, 0x1f170, 0x1f93e, 0x1f004, 0x1f170, 0x1f93f, 0x1f93e}
	testGarbageInput(t, string(runes3[:]), "Emojis from different ecoji versions seen", "mixed_3")

	// validate the assumptions of the test
	if paddingLastV1[0] == paddingLastV2[0] || paddingLastV1[1] == paddingLastV2[1] || paddingLastV1[2] != paddingLastV2[2] || paddingLastV1[3] != paddingLastV2[3] {
		t.Error("Test assumption not valid")
	}

	// the 2nd rune is ecoji v1 only and the 4th runes are ecoji v2 padding
	runes41 := [4]rune{0x1f004, 0x1f170, 0x1f004, paddingLastV2[0]}
	testGarbageInput(t, string(runes41[:]), "Emojis from different ecoji versions seen", "mixed_4")
	runes42 := [4]rune{0x1f004, 0x1f170, 0x1f004, paddingLastV2[1]}
	testGarbageInput(t, string(runes42[:]), "Emojis from different ecoji versions seen", "mixed_5")

	// the 3rd rune is ecoji v2 only and the 4th runes are ecoji v1 padding
	runes51 := [4]rune{0x1f004, 0x1f004, 0x1f93f, paddingLastV1[0]}
	testGarbageInput(t, string(runes51[:]), "Emojis from different ecoji versions seen", "mixed_6")
	runes52 := [4]rune{0x1f004, 0x1f004, 0x1f93f, paddingLastV1[1]}
	testGarbageInput(t, string(runes52[:]), "Emojis from different ecoji versions seen", "mixed_7")
}

type singleRuneWriter struct {
	writer io.Writer
}

func (srw *singleRuneWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		return srw.writer.Write(p[0:1])
	} else {
		return srw.writer.Write(p)
	}
}

func TestSingleByteWriter(t *testing.T) {

	reader := strings.NewReader("👖📸🧈🌭👩☕💲🥇🪚☕")
	buffer1 := &bytes.Buffer{}

	srw := &singleRuneWriter{buffer1}

	err := Decode(reader, srw)
	if err != nil {
		t.Error(err)
	}

	buf, err2 := io.ReadAll(buffer1)
	if err2 != nil {
		t.Error(err2)
	}

	if cmp := bytes.Compare([]byte("abcdefxyz"), buf); cmp != 0 {
		t.Fatalf("single byte writer caused data corruption")
	}
}
