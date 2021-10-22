package ecoji

import (
	"fmt"
	"io"
)

type RuneWriter interface {
	WriteByte(byte) error
	WriteRune(rune) (int, error)
}

func encode(s []byte, w RuneWriter, emojis []rune, padding []rune) (err error) {

	if len(s) == 0 {
		panic("expected data")
	}

	var b0, b1, b2, b3, b4 int = int(s[0]), 0, 0, 0, 0

	if len(s) > 1 {
		b1 = int(s[1])
	}

	if len(s) > 2 {
		b2 = int(s[2])
	}

	if len(s) > 3 {
		b3 = int(s[3])
	}

	if len(s) > 4 {
		b4 = int(s[4])
	}

	runes := []rune{emojis[b0<<2|b1>>6], padding[0], padding[0], padding[0]}

	switch len(s) {
	case 1:
	//nothing to do, all padding
	case 2:
		runes[1] = emojis[(b1&0x3f)<<4|b2>>4]
	case 3:
		runes[1] = emojis[(b1&0x3f)<<4|b2>>4]
		runes[2] = emojis[(b2&0x0f)<<6|b3>>2]
	case 4:
		runes[1] = emojis[(b1&0x3f)<<4|b2>>4]
		runes[2] = emojis[(b2&0x0f)<<6|b3>>2]

		switch b3 & 0x03 {
		case 0:
			runes[3] = padding[1]
		case 1:
			runes[3] = padding[2]
		case 2:
			runes[3] = padding[3]
		case 3:
			runes[3] = padding[4]
		}

	case 5:
		runes[1] = emojis[(b1&0x3f)<<4|b2>>4]
		runes[2] = emojis[(b2&0x0f)<<6|b3>>2]
		runes[3] = emojis[(b3&0x03)<<8|b4]
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
	return encodeAndWrap(r, w, wrap, emojisV1[:], paddingV1[:])
}

//Encodes data using the Ecoji version 2 standard
func EncodeV2(r io.Reader, w RuneWriter, wrap uint) (err error) {
	return encodeAndWrap(r, w, wrap, emojisV2[:], paddingV2[:])
}
