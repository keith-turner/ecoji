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
	if !exists && c != padding && c != padding40 && c != padding41 && c != padding42 && c != padding43 {
		panic("Invalid rune " + string(c))
	}

	return c, s, e
}

//Reads unicode emojis, map each emoji to a 10 bit integer, writes 10 bit intergers
func Decode(r *bufio.Reader, w io.Writer) {
	initMapping()

	for {

		//TODO error check reads
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
		var bits4 int

		switch r4 {
		case padding40:
			bits4 = 0
		case padding41:
			bits4 = 1 << 8
		case padding42:
			bits4 = 2 << 8
		case padding43:
			bits4 = 3 << 8
		default:
			bits4 = revMapping[r4]

		}

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits1 >> 2)
		out[1] = byte((bits1 & 0x3 << 6) | (bits2 >> 4))
		out[2] = byte((bits2 & 0xf << 4) | (bits3 >> 6))
		out[3] = byte((bits3 & 0x3f << 2) | (bits4 >> 8))
		out[4] = byte(bits4 & 0xff)

		switch {
		case r2 == padding:
			out = out[:1]
		case r3 == padding:
			out = out[:2]
		case r4 == padding:
			out = out[:3]
		case r4 == padding40 || r4 == padding41 || r4 == padding42 || r4 == padding43:
			out = out[:4]
		}

		//TODO check err
		w.Write(out)
	}

}
