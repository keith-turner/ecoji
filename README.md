# Ecoji

Ecoji encodes data as emojis.  As a bonus, includes code to decode emojis to original data. 

## Build instructions.

This is my first Go project, so I have no clue what I am doing.  Any tips would be appreciated.

```bash
git clone https://github.com/keith-turner/Ecoji.git
cd Ecoji
export GOPATH=$(pwd)
go install com.github/keith-turner/cmd/decoji/ com.github/keith-turner/cmd/ecoji/
```

## Examples of running

Encode example :

```bash
$  echo "Base64 is so 1999, isn't there something better?" | ./bin/ecoji
🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚
```

Decode example :

```bash
$ echo 🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚 | ./bin/decoji
Base64 is so 1999, isn't there something better?
```

## Technical details

Encoding works by repeatedly reading 10 bits from the input.  Every 10 bit
integer has a unique [Unicode emoji][emoji] character assigned to it.  So for
each 10 bit integer, its assigned emoji is output as utf8.  To decode, this
process is reversed.

[emoji]: https://unicode.org/emoji/
