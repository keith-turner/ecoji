# Ecoji 🏣🔉🦐🔼

**WARNING :** This branch contains a work in progress version of Ecoji V2. Any data encoded using this branch may not be compatible with the final version of Ecoji V2.

Ecoji encodes data as 1024 [emojis][emoji].  It's base1024 with an emoji character set and it preserves sort order.  This repository implements Ecoji using Go.  There are links to other implementations below.  Visit [ecoji.io](https://ecoji.io) to try Ecoji in your browser.

Many have asked how Ecoji compares to base64.  The short answer is that a string encoded with Ecoji will have more bytes, but fewer visible characters, than the same string encoded with base64. With Ecoji, each visible char represents 10 bits, but each character is multi-byte.  With base64 each char represents 6 bits and is one byte.  The following table shows encoding sha256 in different ways.

Encoding | Bytes | Characters 
---------|-------|-----------
none     | 32    | N/A
hex      | 64    | 64
base64   | 44    | 44
ecoji    | 112   | 28

## Installing

Ecoji is published on [snapcraft.io](https://snapcraft.io/ecoji) and can be installed with :

```bash
sudo snap install ecoji
```

## Usage

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

## Examples

### Encoding:

```bash
$ echo "Base64 is so 1999, isn't there something better?" | ecoji
🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑
```

### Decoding:

```bash
$ echo 🏗📩🎦🐇🎛📘🔯🚜💞😽🆖🐊🎱🥁🚄🌱💞😭💮🇵💢🕥🐭🔸🍉🚲🦑🐶💢🕥🔮🔺🍉📸🐮🌼👦🚟🥴📑 | ecoji -d
Base64 is so 1999, isn't there something better?
```

### Concatenation:

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

### Making Hashes More Interesting

```bash
$ cat encode.go  | openssl dgst -binary -sha1 | ecoji
🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉
$ echo 🌰🏐🏡🚟🔶🦅😡😺🚆🍑🕡🦞📍🖊🙀🦉 | ecoji -d | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
$ cat encode.go  | openssl dgst -binary -sha1 | openssl base64
GhAkTyOY/Pta78KImgvofylL19M=
```

(If you want to use Ecoji for hashes, consider the dangers inherent in older systems without utf8 emoji support, different fonts, and similar emojis.)


### A URL Shortener

Four base1024 emojis can represent 1 trillion unique IDs.  In the example below `af82dd48f7` represents a 5 byte id for a URL in a key value store like [Accumulo](https://accumulo.apache.org).  When someone enters the URL, the 5 byte id could be used to obtain the actual URL from the database and then redirect.

```
$ printf "https://ecoji.io/%s\n" $(echo af82dd48f7 | xxd -r -p | ecoji)
https://ecoji.io/😉🈚🛠🏄
```

### Sorting Ecoji-Encoded Data

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

## Implementations

Libraries [implementing](docs/encoding.md) the Ecoji encoding standard. Submit PR to add a library to the table.

| Language | Comments
|----------|----------
| [D](https://github.com/ohdatboi/ecoji-d) | Implementation of Ecoji written in the D programming language.
| Go       | This repository offers a Go library package with two functions [ecoji.Encode()](encode.go) and [ecoji.Decode()](decode.go).
| [Java](https://github.com/netvl/ecoji-java) | Implementation of Ecoji written in Java, usable in any JVM language.
| [.NET](https://github.com/abock/dotnet-ecoji) | Implementation of Ecoji written in C# targeting .NET Standard 2.0: [`dotnet add package Ecoji`](https://www.nuget.org/packages/Ecoji).
| [PHP](https://github.com/Rayne/ecoji-php) | PHP 7.x implementation of Ecoji. Available as [`rayne/ecoji` on Packagist](https://packagist.org/packages/rayne/ecoji).
| [Python](https://github.com/mecforlove/ecoji-py) | Implementation of Ecoji written in the Python3 programming language.
| [Rust](https://github.com/netvl/ecoji.rs) | Implementation of Ecoji written in the Rust programming language.
| [Swift](https://github.com/Robindiddams/ecoji-swift) | Implementation of Ecoji written in the Swift programming language.


## Building

This is my first Go project and I am starting to get my bearings. If you are new
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
