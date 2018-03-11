# Ecoji 🏣🔉🦐🔼

Ecoji encodes data as emojis.  As a bonus, includes code to decode emojis to original data. 

## Build instructions.

This is my first Go project, I am starting to get my bearings. If you are new
to Go I would recommend this [video] and the [tour].

```bash
# The following are general Go setup instructions.  Ignore if you know Go, I am new to it.
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# This will download Ecoji to $GOPATH/src
go get github.com/keith-turner/ecoji

# This will build the ecoji command and put it in $GOPATH/bin
go install github.com/keith-turner/ecoji/cmd/ecoji
```

## Examples of running

Encode example :

```bash
$ echo "Base64 is so 1999, isn't there something better?" | ecoji
🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑
```

Decode example :

```bash
$ echo 🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑 | ecoji -d
Base64 is so 1999, isn't there something better?
```

Concatenation :

```bash
$ echo -n abc | ecoji
👖📸🎈☕
$ echo -n 6789 | ecoji
🎥🤠📠🏍
$ echo XY | ecoji
🐲👡🕟☕
$ echo 👖📸🎈☕🎥🤠📠🏍🐲👡🕟☕ | ecoji -d
abc6789XY
```

Make your hashes more interesting.

```bash
$ cat encode.go  | openssl dgst -binary -sha1 | ecoji
🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉
$ echo 🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉 | ecoji -d | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
$ cat encode.go  | openssl dgst -binary -sha1 | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
```

Data encoded with Ecoji sorts the same as the input data.

```bash
$ echo -n a | ecoji > /tmp/stest.ecoji
$ echo -n ab | ecoji >> /tmp/test.ecoji
$ echo -n abc | ecoji >> /tmp/test.ecoji
$ echo -n abcd | ecoji >> /tmp/test.ecoji
$ echo -n ac | ecoji >> /tmp/test.ecoji
$ echo -n b | ecoji >> /tmp/test.ecoji
$ echo -n ba | ecoji >> /tmp/test.ecoji

$ export LC_ALL=C
$ sort /tmp/test.ecoji > /tmp/test-sorted.ecoji
$ diff /tmp/test.ecoji /tmp/test-sorted.ecoji
$ cat /tmp/test-sorted.ecoji
👖📲☕☕
👖📸🎈☕
👖📸🎦⚜
👖🔃☕☕
👙☕☕☕
👚📢☕☕
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

## Library

Ecoji offers a Go library package with two functions `ecoji.Encode()` and `ecoji.Decode()`.

## Technical details

Encoding works by repeatedly reading 10 bits from the input.  Every 10 bit
integer has a unique [Unicode emoji][emoji] character assigned to it.  So for
each 10 bit integer, its assigned emoji is output as utf8.  To decode, this
process is reversed.

Ecoji is base1024 using a subset of emojis as its numerals.



[emoji]: https://unicode.org/emoji/
[video]: https://www.youtube.com/watch?v=XCsL89YtqCs
[tour]: https://tour.golang.org/welcome/1
