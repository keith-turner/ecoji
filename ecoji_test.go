package ecoji

import (
	"bytes"
	"math/big"
	"math/rand"
	"strings"
	"testing"
)

func checkV1(t *testing.T, expected []rune, input []byte) {
	reader := bytes.NewBuffer(input)
	buffer1 := bytes.NewBuffer(nil)

	err := Encode(reader, buffer1, 0)
	if err != nil {
		t.Error(err)
	}

	actual, _ := buffer1.ReadString('\n')

	if cmp := strings.Compare(actual, string(expected)); cmp != 0 {
		t.Errorf("'%s' != '%s' %d", string(expected), actual, cmp)
	}

	buffer2 := bytes.NewBuffer(nil)
	err = Decode(strings.NewReader(string(expected)), buffer2)
	if err != nil {
		t.Error(err)
	}

	if cmp := bytes.Compare(input, buffer2.Bytes()); cmp != 0 {
		t.Errorf("'%v' != '%v' %d", input, buffer2.Bytes(), cmp)
	}
}

func check(t *testing.T, expected []rune, input []byte) {
	reader := bytes.NewBuffer(input)
	buffer1 := bytes.NewBuffer(nil)

	err := EncodeV2(reader, buffer1, 0)
	if err != nil {
		t.Error(err)
	}

	actual, _ := buffer1.ReadString('\n')

	if cmp := strings.Compare(actual, string(expected)); cmp != 0 {
		t.Errorf("'%s' != '%s' %d", string(expected), actual, cmp)
	}

	buffer2 := bytes.NewBuffer(nil)
	err = Decode(strings.NewReader(string(expected)), buffer2)
	if err != nil {
		t.Error(err)
	}

	if cmp := bytes.Compare(input, buffer2.Bytes()); cmp != 0 {
		t.Errorf("'%v' != '%v' %d", input, buffer2.Bytes(), cmp)
	}
}

func TestOneByteEncode(t *testing.T) {
	check(t, []rune{emojisV2[int('k')<<2], PADDING, PADDING, PADDING}, []byte{'k'})
	checkV1(t, []rune{emojisV1[int('k')<<2], PADDING, PADDING, PADDING}, []byte{'k'})
}

func TestTwoByteEncode(t *testing.T) {
	check(t, []rune{emojisV2[0], emojisV2[16], PADDING, PADDING}, []byte{0x00, 0x01})
	checkV1(t, []rune{emojisV1[0], emojisV1[16], PADDING, PADDING}, []byte{0x00, 0x01})
}

func TestThreeByteEncode(t *testing.T) {
	check(t, []rune{emojisV2[0], emojisV2[16], emojisV2[128], PADDING}, []byte{0x00, 0x01, 0x02})
	checkV1(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], PADDING}, []byte{0x00, 0x01, 0x02})
}

func TestFourByteEncode(t *testing.T) {
	check(t, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[0]}, []byte{0x00, 0x01, 0x02, 0x00})
	check(t, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[1]}, []byte{0x00, 0x01, 0x02, 0x01})
	check(t, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[2]}, []byte{0x00, 0x01, 0x02, 0x02})
	check(t, []rune{emojisV2[0], emojisV2[16], emojisV2[128], paddingLastV2[3]}, []byte{0x00, 0x01, 0x02, 0x03})

	checkV1(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[0]}, []byte{0x00, 0x01, 0x02, 0x00})
	checkV1(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[1]}, []byte{0x00, 0x01, 0x02, 0x01})
	checkV1(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[2]}, []byte{0x00, 0x01, 0x02, 0x02})
	checkV1(t, []rune{emojisV1[0], emojisV1[16], emojisV1[128], paddingLastV1[3]}, []byte{0x00, 0x01, 0x02, 0x03})
}

func TestFiveByteEncode(t *testing.T) {
	check(t, []rune{emojisV2[687], emojisV2[222], emojisV2[960], emojisV2[291]}, []byte{0xab, 0xcd, 0xef, 0x01, 0x23})
}

func TestGarbage(t *testing.T) {
	reader := strings.NewReader("not emojisV2")
	buffer1 := bytes.NewBuffer(nil)

	err := Decode(reader, buffer1)
	if err == nil {
		t.Error("Expected error")
	}

	if !strings.Contains(err.Error(), "Invalid rune") {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestExhaustive(t *testing.T) {
	// use this to hold 10 bit ordinals
	biggy := big.NewInt(1)

	// create an array w/ all 1024 ordinals in random order
	ordinals := rand.Perm(1024)

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

	checkV1(t, expectedRunesV1[:], bytes)
	check(t, expectedRunesV2[:], bytes)
}
