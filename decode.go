package ecoji

import (
	"errors"
	"io"
)

type ecojiver int

const (
	BOTH ecojiver = 1
	V1   ecojiver = 2
	V2   ecojiver = 3
	NONE ecojiver = 4
)

func isPaddingV2(r rune) bool {
	return r == paddingV2[0] || r == paddingV2[1] || r == paddingV2[2] || r == paddingV2[3] || r == paddingV2[4]
}

func isPaddingV1(r rune) bool {
	return r == paddingV1[0] || r == paddingV1[1] || r == paddingV1[2] || r == paddingV1[3] || r == paddingV1[4]
}

func checkRuneV1(r rune) bool {
	if _, exists := revEmojisV1[r]; !exists && !isPaddingV1(r) {
		return false
	}
	return true
}

func checkRuneV2(r rune) bool {
	if _, exists := revEmojisV2[r]; !exists && !isPaddingV2(r) {
		return false
	}
	return true
}

func checkRune(r rune) ecojiver {
	//TODO create a map where the value is ecojiver, will be more efficient
	isV1 := checkRuneV1(r)
	isV2 := checkRuneV2(r)

	if isV1 && isV2 {
		return BOTH
	} else if isV1 {
		return V1
	} else if isV2 {
		return V2
	} else {
		return NONE
	}

}

func readRune(r io.RuneReader, currver *ecojiver) (c rune, size int, err error) {
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
	switch ver := checkRune(c); ver {
	case NONE:
		{
			return 0, 0, errors.New("Invalid rune " + string(c))
		}
	case BOTH:
		{
			//noop
		}
	case V1:
		{
			if *currver == BOTH {
				*currver = V1
			} else if *currver != V1 {
				return 0, 0, errors.New("Invalid rune " + string(c))
			}
		}
	case V2:
		{
			if *currver == BOTH {
				*currver = V2
			} else if *currver != V2 {
				return 0, 0, errors.New("Invalid rune " + string(c))
			}
		}
	default:
		panic("Unexpected ecoji ver")

	}

	return c, s, e
}

//Decodes data encoded using the Ecoji version 1 or 2 standard back to the original data.
func Decode(r io.RuneReader, w io.Writer) (err error) {
	ver := BOTH

	for {
		var runes [4]rune

		//TODO error check reads
		r1, _, e1 := readRune(r, &ver)
		if e1 == io.EOF {
			break
		} else if e1 != nil {
			return e1
		}
		runes[0] = r1

		for i := 1; i < 4; i++ {
			r, _, err := readRune(r, &ver)
			if err == io.EOF {
				return errors.New("Unexpected end of data, input data size not multiple of 4")
			} else if err != nil {
				return err
			}
			runes[i] = r
		}
		var bits1, bits2, bits3, bits4 int

		if ver == V1 {
			bits1 = revEmojisV1[runes[0]]
			bits2 = revEmojisV1[runes[1]]
			bits3 = revEmojisV1[runes[2]]

			switch runes[3] {
			case paddingV1[1]:
				bits4 = 0
			case paddingV1[2]:
				bits4 = 1 << 8
			case paddingV1[3]:
				bits4 = 2 << 8
			case paddingV1[4]:
				bits4 = 3 << 8
			default:
				bits4 = revEmojisV1[runes[3]]
			}
		} else {
			bits1 = revEmojisV2[runes[0]]
			bits2 = revEmojisV2[runes[1]]
			bits3 = revEmojisV2[runes[2]]

			switch runes[3] {
			case paddingV2[1]:
				bits4 = 0
			case paddingV2[2]:
				bits4 = 1 << 8
			case paddingV2[3]:
				bits4 = 2 << 8
			case paddingV2[4]:
				bits4 = 3 << 8
			default:
				bits4 = revEmojisV2[runes[3]]
			}
		}

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits1 >> 2)
		out[1] = byte((bits1 & 0x3 << 6) | (bits2 >> 4))
		out[2] = byte((bits2 & 0xf << 4) | (bits3 >> 6))
		out[3] = byte((bits3 & 0x3f << 2) | (bits4 >> 8))
		out[4] = byte(bits4 & 0xff)

		if ver == V1 {
			switch {
			case runes[1] == paddingV1[0]:
				out = out[:1]
			case runes[2] == paddingV1[0]:
				out = out[:2]
			case runes[3] == paddingV1[0]:
				out = out[:3]
			case runes[3] == paddingV1[1] || runes[3] == paddingV1[2] || runes[3] == paddingV1[3] || runes[3] == paddingV1[4]:
				out = out[:4]
			}
		} else {
			switch {
			case runes[1] == paddingV2[0]:
				out = out[:1]
			case runes[2] == paddingV2[0]:
				out = out[:2]
			case runes[3] == paddingV2[0]:
				out = out[:3]
			case runes[3] == paddingV2[1] || runes[3] == paddingV2[2] || runes[3] == paddingV2[3] || runes[3] == paddingV2[4]:
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
