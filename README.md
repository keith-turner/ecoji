# Ecoji

Ecoji encodes data as emojis.  As a bonus, includes code to decode emojis to original data. 

## Build instructions.

This is my first Go project, so I have no clue what I am doing.  Any tips would be appreciated.

```bash
git clone https://github.com/keith-turner/Ecoji.git
cd Ecoji
export GOPATH=$(pwd)
go install github.com/keith-turner/cmd/ecoji/
```

## Examples of running

Encode example :

```bash
$  echo "Base64 is so 1999, isn't there something better?" | ./bin/ecoji
🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚
```

Decode example :

```bash
$ echo 🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚 | ./bin/ecoji -d
Base64 is so 1999, isn't there something better?
```

Usage :

```bash
$ ./bin/ecoji -h
usage: ecoji [OPTIONS]... [FILE]

Encode or decode data as Unicode emojis. 😁

Options:
    -d, --decode          decode data
    -w, --wrap=COLS       wrap encoded lines after COLS character (default 76).
                          Use 0 to disable line wrapping
```

## Technical details

Encoding works by repeatedly reading 10 bits from the input.  Every 10 bit
integer has a unique [Unicode emoji][emoji] character assigned to it.  So for
each 10 bit integer, its assigned emoji is output as utf8.  To decode, this
process is reversed.

Ecoji is base1024 using a subset of emojis as its numerals.

[emoji]: https://unicode.org/emoji/
