package ecoji

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
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

func check(t *testing.T, expectedV1 []rune, expectedV2 []rune, input []byte, name string) {
	reader := bytes.NewBuffer(input)
	buffer1 := bytes.NewBuffer(nil)

	err := EncodeV2(reader, buffer1, 0)
	if err != nil {
		t.Error(err)
	}

	actual, _ := buffer1.ReadString('\n')

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

	err = Encode(reader2, buffer3, 0)
	if err != nil {
		t.Error(err)
	}

	actual, _ = buffer3.ReadString('\n')

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
		writeBytes(t, input, "test_scripts/data/"+name+".plain")
		writeRunes(t, expectedV1, "test_scripts/data/"+name+".ev1")
		writeRunes(t, expectedV2, "test_scripts/data/"+name+".ev2")
	}
}

func TestOneByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[int('k')<<2], padding, padding, padding}, []rune{emojisV2[int('k')<<2], padding}, []byte{'k'}, "one_byte")
}

func TestTwoByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], padding, padding}, []rune{emojisV2[0], emojisV2[16], padding}, []byte{0x00, 0x01}, "two_byte")
}

func TestThreeByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], padding}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], padding}, []byte{0x00, 0x01, 0x02}, "three_byte")
}

func TestFourByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[0]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[0]}, []byte{0x00, 0x01, 0x02, 0x00}, "four_byte_1")
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[1]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[1]}, []byte{0x00, 0x01, 0x02, 0x01}, "four_byte_2")
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[2]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[2]}, []byte{0x00, 0x01, 0x02, 0x02}, "four_byte_3")
	check(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[3]}, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[3]}, []byte{0x00, 0x01, 0x02, 0x03}, "four_byte_4")
}

func TestFiveByteEncode(t *testing.T) {
	check(t, []rune{emojisV1[687], emojisV1[222], emojisV1[960], emojisV1[291]}, []rune{emojisV2[687], emojisV2[222], emojisV2[960], emojisV2[291]}, []byte{0xab, 0xcd, 0xef, 0x01, 0x23}, "five_byte")
}

func TestSixByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 123<<38 | 456<<28 | 789<<18 | 909<<8 | 55
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[123], emojisV1[456], emojisV1[789], emojisV1[909], emojisV1[55<<2], padding, padding, padding},
		[]rune{emojisV2[123], emojisV2[456], emojisV2[789], emojisV2[909], emojisV2[55<<2], padding}, data[2:8], "six_byte")
}

func TestSevenByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 237<<46 | 77<<36 | 257<<26 | 513<<16 | 809<<6 | 7
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[237], emojisV1[77], emojisV1[257], emojisV1[513], emojisV1[809], emojisV1[7<<4], padding, padding},
		[]rune{emojisV2[237], emojisV2[77], emojisV2[257], emojisV2[513], emojisV2[809], emojisV2[7<<4], padding}, data[1:8], "seven_byte")
}

func TestEightByteEncode(t *testing.T) {
	var scratch uint64
	scratch = 3<<54 | 206<<44 | 368<<34 | 617<<24 | 650<<14 | 1005<<4 | 3
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], scratch)

	check(t, []rune{emojisV1[3], emojisV1[206], emojisV1[368], emojisV1[617], emojisV1[650], emojisV1[1005], emojisV1[3<<6], padding},
		[]rune{emojisV2[3], emojisV2[206], emojisV2[368], emojisV2[617], emojisV2[650], emojisV2[1005], emojisV1[3<<6], padding}, data[0:8], "eight_byte")
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
			[]rune{emojisV2[855], emojisV2[298], emojisV2[1007], emojisV2[97], emojisV2[611], emojisV2[291], emojisV2[856], paddingLastV2[i]}, data[0:9], "nine_byte_"+fmt.Sprint(i))
	}
}

func TestGarbage(t *testing.T) {
	reader := strings.NewReader("not emojisV2")
	buffer1 := bytes.NewBuffer(nil)

	err := Decode(reader, buffer1)
	if err == nil {
		t.Error("Expected error")
	}

	if !strings.Contains(err.Error(), "Non Ecoji character seen") {
		t.Errorf("Unexpected error message: %s", err.Error())
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

	check(t, expectedRunesV1[:], expectedRunesV2[:], bytes, "exhaustive")
}

func TestPhrase(t *testing.T) {
	expectedV1 := []rune("🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑")
	expectedV2 := []rune("🧗📩🧊🐇🧇📘🔯🚜💞😽♑🐊🎱🥁🚄🌱💞😭💮✨💢🪠🐭🩴🍉🚲🦑🐶💢🪠🔮🩹🍉📸🐮🌼👦🚟🥴📑")
	plain := []byte("Base64 is so 1999, isn't there something better?\n")
	check(t, expectedV1, expectedV2, plain, "phrase")
}
