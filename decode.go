package ecoji

import (
	"errors"
	"io"
)

func isPadding(r rune) bool {
	return r == padding || r == padding40 || r == padding41 || r == padding42 || r == padding43
}

func isPaddingV1(r rune) bool {
	return r == paddingV1 || r == padding40V1 || r == padding41V1 || r == padding42V1 || r == padding43V1
}

func checkRuneV1(r rune) bool {
	if _, exists := revEmojisV1[r]; !exists && !isPaddingV1(r) {
		return false
	}
	return true
}

func checkRune(r rune) bool {
	if _, exists := revEmojis[r]; !exists && !isPadding(r) {
		return false
	}
	return true
}

func readRune(r io.RuneReader, ecojiV1 *bool) (c rune, size int, err error) {
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
	if *ecojiV1 {
		if ok := checkRuneV1(c); !ok {
			return 0, 0, errors.New("Invalid rune " + string(c))
		}
	} else {
		if ok := checkRune(c); !ok {
			// try to fallback to ecoji v1
			if isV1 := checkRuneV1(c); isV1 {
				*ecojiV1 = true
			} else {
				return 0, 0, errors.New("Invalid rune " + string(c))
			}
		}
	}
	return c, s, e
}

//Reads unicode emojis, map each emoji to a 10 bit integer, writes 10 bit intergers
func Decode(r io.RuneReader, w io.Writer) (err error) {
	var ecojiV1 bool

	for {
		var runes [4]rune

		//TODO error check reads
		r1, _, e1 := readRune(r, &ecojiV1)
		if e1 == io.EOF {
			break
		} else if e1 != nil {
			return e1
		}
		runes[0] = r1

		for i := 1; i < 4; i++ {
			r, _, err := readRune(r, &ecojiV1)
			if err == io.EOF {
				return errors.New("Unexpected end of data, input data size not multiple of 4")
			} else if err != nil {
				return err
			}
			runes[i] = r
		}
		var bits1, bits2, bits3, bits4 int

		if ecojiV1 {
			bits1 = revEmojisV1[runes[0]]
			bits2 = revEmojisV1[runes[1]]
			bits3 = revEmojisV1[runes[2]]

			switch runes[3] {
			case padding40V1:
				bits4 = 0
			case padding41V1:
				bits4 = 1 << 8
			case padding42V1:
				bits4 = 2 << 8
			case padding43V1:
				bits4 = 3 << 8
			default:
				bits4 = revEmojisV1[runes[3]]
			}
		} else {
			bits1 = revEmojis[runes[0]]
			bits2 = revEmojis[runes[1]]
			bits3 = revEmojis[runes[2]]

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
		}

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits1 >> 2)
		out[1] = byte((bits1 & 0x3 << 6) | (bits2 >> 4))
		out[2] = byte((bits2 & 0xf << 4) | (bits3 >> 6))
		out[3] = byte((bits3 & 0x3f << 2) | (bits4 >> 8))
		out[4] = byte(bits4 & 0xff)

		if ecojiV1 {
			switch {
			case runes[1] == paddingV1:
				out = out[:1]
			case runes[2] == paddingV1:
				out = out[:2]
			case runes[3] == paddingV1:
				out = out[:3]
			case runes[3] == padding40V1 || runes[3] == padding41V1 || runes[3] == padding42V1 || runes[3] == padding43V1:
				out = out[:4]
			}
		} else {
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
