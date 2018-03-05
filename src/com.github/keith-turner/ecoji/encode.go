package ecoji

import (
	"bufio"
	"io"
)

func encode(s []byte, w *bufio.Writer) {

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

	// read 8 bits from 1st byte and 2 bits from 2nd byte
	w.WriteRune(mapping[b0<<2|b1>>6])

	if len(s) >= 2 {
		// read 6 bits from 2nd byte and 4 bits from 3rd byte
		w.WriteRune(mapping[(b1&0x3f)<<4|b2>>4])
	}
	if len(s) >= 3 {
		// read 4 bits from 3rd byte and 6 bits from 4th byte
		w.WriteRune(mapping[(b2&0x0f)<<6|b3>>2])
	}
	if len(s) >= 4 {
		// read 2 bits from 4th byte and 8 bits from 5th byte
		w.WriteRune(mapping[(b3&0x03)<<8|b4])
	}

	if len(s) < 5 {
		w.WriteRune(endRune)
	}
}

func readFully(r io.Reader, buffer []byte) (n int, e error) {
	num, err := r.Read(buffer)

	for num < len(buffer) && err != io.EOF {
		more, err2 := r.Read(buffer[num:])
		num += more
		err = err2
	}

	return num, err
}

//Maps every 10 bits from the reader to one of 1024 Unicode emojis, writing the emojis.
func Encode(r io.Reader, w *bufio.Writer) {

	initMapping()

	buffer := make([]byte, 5)
	printed := 0

	for {
		num, err := readFully(r, buffer)

		if num == 0 && err == io.EOF {
			if printed > 0 {
				w.WriteByte('\n')
			}
			break
		}
		//TODO check err
		encode(buffer[0:num], w)
		printed += 4
		if printed >= 72 {
			w.WriteByte('\n')
			printed = 0
		}

	}
}
