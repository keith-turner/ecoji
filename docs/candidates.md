# Ecoji candidates

This file [candidates.txt] contains all of the emojis that were candidates for
Ecoji 2.  This was created by downloading [emoji-test.txt][1] and running the
following command on it.

```bash
grep -v "^#" emoji-test.txt | grep -v -e '^[[:space:]]*$' | grep -E '^[0-9A-F]+[[:space:]]+;' | grep "fully-qualified" | cut -f 1 -d ' '> docs/candidates.txt
```

The purpose of this command is to find all fully qualified single code point emojis.


[1]: https://unicode.org/Public/emoji/13.0/emoji-test.txt
