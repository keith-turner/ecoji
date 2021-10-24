package ecoji

import (
	"errors"
	"fmt"
	"io"
)

func nextRune(r io.RuneReader, expectedVer ecojiver) (emojiInfo, ecojiver, error) {
	c, _, e := r.ReadRune()

	if e != nil {
		return emojiInfo{}, -1, e
	}

	for c == '\n' {
		c, _, e = r.ReadRune()
		if e != nil {
			return emojiInfo{}, -1, e
		}
	}

	einfo, exists := revEmojis[c]

	if !exists {
		return emojiInfo{}, -1, errors.New("Invalid rune " + string(c))
	}

	if einfo.version != EVALL {
		if expectedVer == EVALL {
			expectedVer = einfo.version
		} else if expectedVer != einfo.version {
			return emojiInfo{}, -1, errors.New("Emojis from different ecoji versions seen " + string(c))
		}
	}

	return einfo, expectedVer, nil
}

//Decodes data encoded using the Ecoji version 1 or 2 standard back to the original data.
func Decode(r io.RuneReader, w io.Writer) error {
	expectedVer := EVALL

	for {
		var emojis [4]emojiInfo

		for i := 0; i < 4; i++ {
			var err error
			var ei emojiInfo
			ei, expectedVer, err = nextRune(r, expectedVer)
			if err == io.EOF {
				if i == 0 {
					return nil
				} else {
					return errors.New("Unexpected end of data, input data size not multiple of 4")
				}
			} else if err != nil {
				return err
			}
			emojis[i] = ei
		}

		paddingIsValid := emojis[0].padding == PAD_NONE && ((emojis[1].padding == PAD_NONE && emojis[2].padding == PAD_NONE && emojis[3].padding == PAD_NONE) ||
			(emojis[1].padding == PAD_NONE && emojis[2].padding == PAD_NONE && (emojis[3].padding == PAD_FILL || emojis[3].padding == PAD_LAST)) ||
			(emojis[1].padding == PAD_NONE && emojis[2].padding == PAD_FILL && emojis[3].padding == PAD_FILL) ||
			(emojis[1].padding == PAD_FILL && emojis[2].padding == PAD_FILL && emojis[3].padding == PAD_FILL))

		if !paddingIsValid {
			return fmt.Errorf("Unexpected padding seen %v", emojis)
		}

		bits := int64(emojis[0].ordinal)<<30 |
			int64(emojis[1].ordinal)<<20 |
			int64(emojis[2].ordinal)<<10 |
			int64(emojis[3].ordinal)

		out := []byte{0, 0, 0, 0, 0}

		out[0] = byte(bits >> 32)
		out[1] = byte(bits >> 24)
		out[2] = byte(bits >> 16)
		out[3] = byte(bits >> 8)
		out[4] = byte(bits)

		switch {
		case emojis[1].padding == PAD_FILL:
			out = out[:1]
		case emojis[2].padding == PAD_FILL:
			out = out[:2]
		case emojis[3].padding == PAD_FILL:
			out = out[:3]
		case emojis[3].padding == PAD_LAST:
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
}
