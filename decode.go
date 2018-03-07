package ecoji

import (
	"bufio"
	"io"
)

func ReadRune(r *bufio.Reader) (c rune, size int, err error) {
	c, s, e := r.ReadRune()
	for c == '\n' {
		c, s, e = r.ReadRune()
	}

	if e == io.EOF {
		return c, s, e
	}

	// check to see if this is a valid emoji rune
	_, exists := revMapping[c]
	if !exists && c != endRune {
		panic("Invalid rune " + string(c))
	}

	return c, s, e
}

//Reads unicode emojis, map each emoji to a 10 bit integer, writes 10 bit intergers
func Decode(r *bufio.Reader, w io.Writer) {
	initMapping()

	for {

		r1, _, e1 := ReadRune(r)
		if e1 == io.EOF {
			break
		}

		r2, _, _ := ReadRune(r)
		r3, _, _ := ReadRune(r)
		r4, _, _ := ReadRune(r)

		bits1 := revMapping[r1]
		bits2 := revMapping[r2]
		bits3 := revMapping[r3]
		bits4 := revMapping[r4]

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits1 >> 2)
		out[1] = byte((bits1 & 0x3 << 6) | (bits2 >> 4))
		out[2] = byte((bits2 & 0xf << 4) | (bits3 >> 6))
		out[3] = byte((bits3 & 0x3f << 2) | (bits4 >> 8))
		out[4] = byte(bits4 & 0xff)

		if r2 == endRune {
			out = out[:1]
		} else if r3 == endRune {
			out = out[:2]
		} else if r4 == endRune {
			out = out[:3]
		} else {
			r5, _, e5 := ReadRune(r)
			if e5 != io.EOF {
				if r5 == endRune {
					out = out[:4]
				} else {
					r.UnreadRune()
					//TODO check err
				}
			}
		}

		//TODO check err
		w.Write(out)
	}

}
