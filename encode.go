package ecoji

import (
	"fmt"
	"io"
)

type RuneWriter interface {
	WriteByte(byte) error
	WriteRune(rune) (int, error)
}

func encode(s []byte, w RuneWriter, emojis []rune, paddingLast []rune) error {

	if len(s) == 0 {
		panic("expected data")
	}

	var bits uint64

	switch len(s) {
	case 1:
		bits = uint64(s[0]) << 32
	case 2:
		bits = uint64(s[0])<<32 | uint64(s[1])<<24
	case 3:
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16
	case 4:
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16 | uint64(s[3])<<8
	case 5:
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16 | uint64(s[3])<<8 | uint64(s[4])
	}

	runes := []rune{emojis[bits>>30], PADDING, PADDING, PADDING}

	switch len(s) {
	case 1:
	//nothing to do, all padding
	case 2:
		runes[1] = emojis[0x3ff&(bits>>20)]
	case 3:
		runes[1] = emojis[0x3ff&(bits>>20)]
		runes[2] = emojis[0x3ff&(bits>>10)]
	case 4:
		runes[1] = emojis[0x3ff&(bits>>20)]
		runes[2] = emojis[0x3ff&(bits>>10)]
		runes[3] = paddingLast[(0x03 & (bits >> 8))]
	case 5:
		runes[1] = emojis[0x3ff&(bits>>20)]
		runes[2] = emojis[0x3ff&(bits>>10)]
		runes[3] = emojis[0x3ff&bits]
	default:
		panic(fmt.Sprintf("unexpected length %d", len(s)))

	}

	for _, r := range runes {
		if _, err := w.WriteRune(r); err != nil {
			return err
		}
	}

	return nil
}

func readFully(r io.Reader, buffer []byte) (n int, e error) {
	num, err := r.Read(buffer)

	for num < len(buffer) && err != io.EOF && err == nil {
		more, err2 := r.Read(buffer[num:])
		num += more
		err = err2
	}

	return num, err
}

func encodeAndWrap(r io.Reader, w RuneWriter, wrap uint, emojis []rune, padding []rune) (err error) {
	buffer := make([]byte, 5)
	printed := uint(0)

	for {

		num, err := readFully(r, buffer)

		if num == 0 && err == io.EOF {
			if printed > 0 {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			break
		}

		if err != nil && err != io.EOF {
			return err
		}

		if err := encode(buffer[0:num], w, emojis, padding); err != nil {
			return err
		}

		if wrap > 0 {
			printed += 4
			if printed >= wrap {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
				printed = 0
			}
		}

	}

	return nil
}

//Encodes data using the Ecoji version 1 standard
func Encode(r io.Reader, w RuneWriter, wrap uint) (err error) {
	return encodeAndWrap(r, w, wrap, emojisV1[:], paddingLastV1[:])
}

//Encodes data using the Ecoji version 2 standard
func EncodeV2(r io.Reader, w RuneWriter, wrap uint) (err error) {
	return encodeAndWrap(r, w, wrap, emojisV2[:], paddingLastV2[:])
}
