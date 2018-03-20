# Ecoji 🏣🔉🦐🔼

Ecoji encodes data as 1024 [emojis][emoji], its base1024 with an emoji character set.  As a bonus, includes code to decode emojis to original data.  

Many have asked how Ecoji compares to base64.  The short answer is that Ecoji is more bytes, but less visible characters. For example sha256 is 112 bytes when encoded using Ecoji. For base64, sha256 is 44 bytes. For sha256, Ecoji has 28 visible chars and base64 44. With base64 each visible char encodes 6 bits of data and each visible char is one byte. With Ecoji each visible char encodes 10 bits of data, but each visible character is multi-byte.

## Installing

Ecoji is published to [snapcraft.io](https://snapcraft.io/ecoji) and can be installed with :

```bash
sudo snap install ecoji
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

Make your hashes more interesting. When using for hashes, consideration should be given for older systems without utf8 emoji support, different fonts, and similar emojis.

```bash
$ cat encode.go  | openssl dgst -binary -sha1 | ecoji
🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉
$ echo 🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉 | ecoji -d | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
$ cat encode.go  | openssl dgst -binary -sha1 | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
```

Make a cool URL shortener.  Four base1024 emojis can represent 1 trillion unique IDs.  In the example below `af82dd48f7` represents a 5 byte id for a URL in a key value store like [Accumulo](https://accumulo.apache.org).  When someone enters the URL, the 5 byte id could be used to obtain the actual URL from the database and then redirect.

```
$ printf "https://ecoji.io/%s\n" $(echo af82dd48f7 | xxd -r -p | ecoji)
https://ecoji.io/😉🈚🛠🏄
```

Data encoded with Ecoji sorts the same as the input data.

```bash
$ echo -n a | ecoji > /tmp/test.ecoji
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
👕☕☕☕
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
    -h, --help            Print this message
    -v, --version         Print version information.
```

## Libraries

Libraries [implementing](docs/encoding.md) the Ecoji encoding standard. Submit PR to add a library to the table.

| Language | Comments
|----------|----------
| [D](https://github.com/ohdatboi/ecoji-d) | Implementation of Ecoji written in the D programming language.
| Go       | This repository offers a Go library package with two functions [ecoji.Encode()](encode.go) and [ecoji.Decode()](decode.go).
| Java     | Coming soon, I plan to implement this and publish to maven central unless someone else does.
| [PHP](https://github.com/Rayne/ecoji-php) | PHP 7.x implementation of Ecoji. Available as [`rayne/ecoji` on Packagist](https://packagist.org/packages/rayne/ecoji).
| [Rust](https://github.com/netvl/ecoji.rs) | Implementation of Ecoji written in the Rust programming language.

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

[emoji]: https://unicode.org/emoji/
[video]: https://www.youtube.com/watch?v=XCsL89YtqCs
[tour]: https://tour.golang.org/welcome/1
