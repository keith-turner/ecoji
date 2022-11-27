package ecoji

import (
	"fmt"
	"io"
)

type RuneWriter interface {
	WriteByte(byte) error
	WriteRune(rune) (int, error)
}

func encode(s []byte, w RuneWriter, emojis []rune, paddingLast []rune, trim bool) error {

	if len(s) == 0 {
		// this should not happen unless there is a bug in this code, hence the panic
		panic("expected data")
	}

	// This variable is used as temp space for conversion. The code will shift input
	// bytes into this variable and then shift 10 bits out at time to map to an
	// emoji. The code below uses 0x3ff a lot which is hex, in binary this is
	// 1111111111 which is 10 1 bits. A mask like (0x3ff & bits) gets the last 10
	// bits from the 64-bit integer.
	var bits uint64

	var runes []rune

	switch len(s) {
	case 1:
		// There is only a single byte of input so convert it to a 10 bit integer and
		// lookup the emoji for encoding. The padding emoji at postion 2 indicates the
		// input was 1 byte or 8 bits. Optionally pad positions 3 and 4.
		if trim {
			runes = []rune{emojis[uint64(s[0])<<2], padding}
		} else {
			runes = []rune{emojis[uint64(s[0])<<2], padding, padding, padding}
		}
	case 2:
		// Shift 2 bytes for a total of 16 bits into the temp var.
		bits = uint64(s[0])<<32 | uint64(s[1])<<24
		// Extract 2 10 bit integers and use them to lookup 2 emojis for encoding. Using
		// padding for the 3rd emoji (and optionally the 4th). Only 6 bits are set in the
		// last 10 bit integer and that is ok because padding in the 3rd position
		// indicates the input was 2 bytes or 16 bits.
		if trim {
			runes = []rune{emojis[bits>>30], emojis[0x3ff&(bits>>20)], padding}
		} else {
			runes = []rune{emojis[bits>>30], emojis[0x3ff&(bits>>20)], padding, padding}
		}
	case 3:
		// Shift 3 bytes for a total of 24 bits into the temp var.
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16
		// Extract 3 10 bit integers and use them to lookup 3 emojis for encoding. Use
		// padding for the last emoji. Only 4 bits are set in the last 10 bit integers
		// and that is ok because padding emoji in the 4th position indicates the input
		// was 3 bytes or 24 bits.
		runes = []rune{emojis[bits>>30], emojis[0x3ff&(bits>>20)], emojis[0x3ff&(bits>>10)], padding}
	case 4:
		// Shift 4 bytes for a total of 32 bits into the temp var
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16 | uint64(s[3])<<8
		// Since there are 32 bits, extract 3 10 bit integers leaving 2 bits. The 3 10
		// bit integers are used to lookup 3 emojis for encoding. Then use the last 2
		// bits to lookup special padding emojis that encode the 2 bits.
		runes = []rune{emojis[bits>>30], emojis[0x3ff&(bits>>20)], emojis[0x3ff&(bits>>10)], paddingLast[(0x03 & (bits >> 8))]}
	case 5:
		// Shift 5 bytes for a total of 40 bits into the temp var.
		bits = uint64(s[0])<<32 | uint64(s[1])<<24 | uint64(s[2])<<16 | uint64(s[3])<<8 | uint64(s[4])
		// Extract 4 10 bit integers and use them to lookup 4 emojis for encoding.
		runes = []rune{emojis[bits>>30], emojis[0x3ff&(bits>>20)], emojis[0x3ff&(bits>>10)], emojis[0x3ff&bits]}
	default:
		// this should not happen unless there is a bug in this code, hence the panic
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

type wrappingWriter struct {
	writer  RuneWriter
	printed uint
	wrap    uint
}

func (ww *wrappingWriter) WriteByte(b byte) error {
	return ww.writer.WriteByte(b)
}
func (ww *wrappingWriter) WriteRune(r rune) (int, error) {
	num, err := ww.writer.WriteRune(r)
	if err == nil {
		ww.printed++
		if ww.printed%ww.wrap == 0 {
			err := ww.writer.WriteByte('\n')
			if err != nil {
				return 0, err
			}
		}
	}

	return num, err
}

func (ww *wrappingWriter) wrapEnd() error {
	if ww.printed%ww.wrap != 0 {
		err := ww.writer.WriteByte('\n')
		if err != nil {
			return err
		}
	}

	return nil
}

func encodeAndWrap(r io.Reader, w RuneWriter, wrap uint, emojis []rune, padding []rune, trim bool) (err error) {
	buffer := make([]byte, 5)

	var writer RuneWriter
	var ww *wrappingWriter

	if wrap > 0 {
		ww = &wrappingWriter{w, 0, wrap}
		writer = ww
	} else {
		ww = nil
		writer = w
	}

	for {
		num, err := readFully(r, buffer)

		if num == 0 && err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			return err
		}

		if err := encode(buffer[0:num], writer, emojis, padding, trim); err != nil {
			return err
		}
	}

	if ww != nil {
		if err2 := ww.wrapEnd(); err2 != nil {
			return err2
		}
	}

	return nil
}

//Encodes data using the Ecoji version 1 standard
func Encode(r io.Reader, w RuneWriter, wrap uint) (err error) {
	return encodeAndWrap(r, w, wrap, emojisV1[:], paddingLastV1[:], false)
}

//Encodes data using the Ecoji version 2 standard
func EncodeV2(r io.Reader, w RuneWriter, wrap uint) (err error) {
	return encodeAndWrap(r, w, wrap, emojisV2[:], paddingLastV2[:], true)
}
