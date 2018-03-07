# Ecoji

Ecoji encodes data as emojis.  As a bonus, includes code to decode emojis to original data. 

## Build instructions.

This is my first Go project, I am starting to get my bearings. If you are new
to Go I would recommend this [video] and the [tour].

```bash
# The following are general go setup instructions.  Ignore if you know Go, I am new to it.
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# This will download project to $GOPATH/src
go get github.com/keith-turner/ecoji

# This will build the ecoji command and put it in $GOPATH/bin
go install github.com/keith-turner/ecoji/cmd/ecoji
```

## Examples of running

Encode example :

```bash
$  echo "Base64 is so 1999, isn't there something better?" | ecoji
🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚
```

Decode example :

```bash
$ echo 🏖📧🎦🐆🎛📖🔭🚙💝😻🆖🐉🎱🤽🚁🌱💝😫💭🇵💡🕣🐬🔶🍉🚯🦎🐵💡🕣🔬🔸🍉📶🐭🌼👥🚜🥯📐🔚 | ecoji -d
Base64 is so 1999, isn't there something better?
```

Usage :

```bash
$ ecoji -h
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

## Library

I am new to Go, but I think this project also builds a library that anyone could use.

[emoji]: https://unicode.org/emoji/
[video]: https://www.youtube.com/watch?v=XCsL89YtqCs
[tour]: https://tour.golang.org/welcome/1
