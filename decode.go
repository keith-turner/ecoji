package ecoji

import (
	"errors"
	"fmt"
	"io"
)

func readFour(r io.RuneReader, expectedVer *ecojiver, emojis []emojiInfo) (int, error) {

	index := 0
	sawPadding := false

	for index < 4 {
		c, _, e := r.ReadRune()

		if e == io.EOF {
			if index == 0 {
				return 0, nil
			} else if sawPadding && (*expectedVer == evAll || *expectedVer == ev2) {
				// Ecoji V2 trims padding and does not pad out to all 4 chars. Therefore, the
				// last bit of data may not be four runes. Let's go ahead and fill the remaining
				// slots w/ padding to make decoding easier.
				for ; index < 4; index++ {
					emojis[index] = revEmojis[padding]
				}
				return index, nil
			} else {
				// Expect Ecoji V1 to always pad out to 4 runes.
				return -1, errors.New("Unexpected end of data, input data size not multiple of 4")
			}
		} else if e != nil {
			return -1, e
		}

		if c == '\n' {
			continue
		}

		einfo, exists := revEmojis[c]

		if !exists {
			return -1, errors.New("Non Ecoji character seen : " + string(c))
		}

		// If we ever see an emoji that is only used by Ecoji V1 then after that point we
		// should never expect to see an emoji that is only used by Ecoji V2. The inverse
		// is also true, if we ever see an emoji only used by Ecoji V2 then we would not
		// expect to see an emoji only used by Ecoji V1 after that. The following looks
		// for these types of malformed input.
		if einfo.version != evAll {
			if *expectedVer == evAll {
				*expectedVer = einfo.version
			} else if *expectedVer != einfo.version {
				return -1, errors.New("Emojis from different ecoji versions seen " + string(c))
			}
		}

		switch einfo.padding {
		case padNone:
			{
				if sawPadding {
					if *expectedVer == evAll || *expectedVer == ev2 {
						// For Ecoji V2 it may trim padding and not pad out all 4 chars. So this could be
						// concatenated Ecoji data, therefore let's put the rune back and return so the
						// data up to the padding can be decoded
						rs, ok := r.(io.RuneScanner)
						if !ok {
							return -1, errors.New("Unable to handle concatenated data because could not cast to RuneScanner")
						}
						rs.UnreadRune()
						for ; index < 4; index++ {
							emojis[index] = revEmojis[padding]
						}
						return index, nil
					} else {
						// Ecoji V1 would always pad out to 4 runes.  So if concatenating Ecoji v1 data we would expect
						// to see non-padding here
						return -1, errors.New("Unexpectedly saw non-padding after padding")
					}
				}
			}
		case padFill:
			{
				if index == 0 {
					return -1, fmt.Errorf("Padding unexpectedly seen in first position %s", string(c))
				}
				sawPadding = true
			}
		case padLast:
			{
				if index != 3 {
					return -1, fmt.Errorf("Last padding seen in unexpected position %s", string(c))
				}
			}
		}

		emojis[index] = einfo
		index = index + 1
	}

	return index, nil
}

//Decodes data encoded using the Ecoji version 1 or 2 standard back to the original data.
func Decode(r io.RuneReader, w io.Writer) error {
	expectedVer := evAll

	for {
		var emojis [4]emojiInfo

		numRead, err := readFour(r, &expectedVer, emojis[:])

		if err != nil {
			return err
		}
		if numRead == 0 {
			return nil
		}

		// Take the 4 10 bit ordinals associated with each emojis and shift them all into a 64 bit integer.
		bits := int64(emojis[0].ordinal)<<30 |
			int64(emojis[1].ordinal)<<20 |
			int64(emojis[2].ordinal)<<10 |
			int64(emojis[3].ordinal)

		out := []byte{0, 0, 0, 0, 0}

		// Shift 5 8-bit integers out of the 40 bit integer created above to recover the encoded data.
		out[0] = byte(bits >> 32)
		out[1] = byte(bits >> 24)
		out[2] = byte(bits >> 16)
		out[3] = byte(bits >> 8)
		out[4] = byte(bits)

		// Inspect the input data for padding emojis to determine the length of the encoded data.
		switch {
		case emojis[1].padding == padFill:
			out = out[:1]
		case emojis[2].padding == padFill:
			out = out[:2]
		case emojis[3].padding == padFill:
			out = out[:3]
		case emojis[3].padding == padLast:
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
