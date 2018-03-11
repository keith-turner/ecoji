package ecoji

import (
	"errors"
	"io"
)

type RuneReader interface {
	ReadRune() (rune, int, error)
}

func readRune(r RuneReader) (c rune, size int, err error) {
	c, s, e := r.ReadRune()
	for c == '\n' {
		c, s, e = r.ReadRune()
	}

	if e == io.EOF {
		return c, s, e
	} else if e != nil {
		return c, s, e
	}

	// check to see if this is a valid emoji rune
	_, exists := revEmojis[c]
	if !exists && c != padding && c != padding40 && c != padding41 && c != padding42 && c != padding43 {
		return 0, 0, errors.New("Invalid rune " + string(c))
	}

	return c, s, e
}

//Reads unicode emojis, map each emoji to a 10 bit integer, writes 10 bit intergers
func Decode(r RuneReader, w io.Writer) (err error) {

	for {
		var runes [4]rune

		//TODO error check reads
		r1, _, e1 := readRune(r)
		if e1 == io.EOF {
			break
		} else if e1 != nil {
			return e1
		}
		runes[0] = r1

		for i := 1; i < 4; i++ {
			r, _, err := readRune(r)
			if err == io.EOF {
				return errors.New("Unexpected end of data, input data size not multiple of 4")
			} else if err != nil {
				return err
			}
			runes[i] = r
		}

		bits1 := revEmojis[runes[0]]
		bits2 := revEmojis[runes[1]]
		bits3 := revEmojis[runes[2]]
		var bits4 int

		switch runes[3] {
		case padding40:
			bits4 = 0
		case padding41:
			bits4 = 1 << 8
		case padding42:
			bits4 = 2 << 8
		case padding43:
			bits4 = 3 << 8
		default:
			bits4 = revEmojis[runes[3]]

		}

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits1 >> 2)
		out[1] = byte((bits1 & 0x3 << 6) | (bits2 >> 4))
		out[2] = byte((bits2 & 0xf << 4) | (bits3 >> 6))
		out[3] = byte((bits3 & 0x3f << 2) | (bits4 >> 8))
		out[4] = byte(bits4 & 0xff)

		switch {
		case runes[1] == padding:
			out = out[:1]
		case runes[2] == padding:
			out = out[:2]
		case runes[3] == padding:
			out = out[:3]
		case runes[3] == padding40 || runes[3] == padding41 || runes[3] == padding42 || runes[3] == padding43:
			out = out[:4]
		}

		for {
			num, err := w.Write(out)

			if err != nil {
				return err
			}

			if num == len(out) {
				break
			}

			out = out[num:]
		}
	}

	return nil
}
