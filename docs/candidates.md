# Ecoji candidates

This file [candidates.txt](candidates.txt) contains all of the emojis that were candidates for
Ecoji 2.  This was created by downloading [emoji-test.txt][1] and running the
following.

```bash
wget https://unicode.org/Public/emoji/14.0/emoji-test.txt
grep -v "^#" emoji-test.txt | grep -v -e '^[[:space:]]*$' | grep -E '^[0-9A-F]+[[:space:]]+;' | grep "fully-qualified" | sed -e 's/[[:space:]]*; fully-qualified[[:space:]]*/ /' > candidates.txt
```

The purpose of this command is to find all fully qualified single code point emojis.

[1]: https://unicode.org/Public/emoji/14.0/emoji-test.txt
