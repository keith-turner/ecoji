# Ecoji 🏣🔉🦐🔼

**WARNING :** This branch contains a work in progress version of Ecoji V2. Any data encoded using this branch may not be compatible with the final version of Ecoji V2.

Ecoji encodes data using 1024 [emojis][emoji]. This repository contains the canonical implementation of the [Ecoji standard](docs/encoding.md) using Go. Version 2 of the Ecoji standard was release in 2021 with an improved set of emojis.  Visit [ecoji.io](https://ecoji.io) to try Ecoji in your browser.
## Installing

| Platform | Install |
| -------- | ------- |
| Windows | TODO |
| Mac | TODO |
| Ubuntu | Run : `sudo snap install ecoji` |

## Usage

```bash
$ ecoji -h
usage: ecoji [OPTIONS]... [FILE]

Encode or decode data as Unicode emojis. 😁

For compatability, when given no options stdin will be encoded using Ecoji 
version 1. When using the new -e option, stdin is encoded using Ecoji 
version 2.  The -e and -d options are mutually exclusive.

Options:
    -e, --encode          Encode data using Ecoji version 2.  Omitting this
                          option will encode using Ecoji version 1.
    -d, --decode          Decodes data encoded using the Ecoji version 1 or 2 standard.
    -w, --wrap=COLS       wrap encoded lines after COLS character (default 76).
                          Use 0 to disable line wrapping.  This options is
                          ignored when decoding.
    -h, --help            Print this message
    -v, --version         Print version information.
```

## Examples

### Encoding:

```bash
$ echo "Base64 is so 1999, isn't there something better?" | ecoji -e
🧛📩🎦🐇🧋📘🔯🚜💞😽♑🐊🎱🥁🚄🌱💞😭💮✨💢🕥🐭🧿🍉🚲🦑🐶💢🕥🔮🩱🍉📸🐮🌼👦🚟🥴📑
```

### Decoding:

```bash
$ echo 🧛📩🎦🐇🧋📘🔯🚜💞😽♑🐊🎱🥁🚄🌱💞😭💮✨💢🕥🐭🧿🍉🚲🦑🐶💢🕥🔮🩱🍉📸🐮🌼👦🚟🥴📑 | ecoji -d
Base64 is so 1999, isn't there something better?
```

### Concatenation:

```bash
$ echo -n abc | ecoji -e
👖📸🎈☕
$ echo -n 6789 | ecoji -e
🎥🤠📠🛼
$ echo XY | ecoji -e
🐲👡🕟☕
$ echo 👖📸🎈☕🎥🤠📠🛼🐲👡🕟☕ | ecoji -d
abc6789XY
```

### Making Hashes More Interesting

```bash
$ cat encode.go  | openssl dgst -binary -sha1 | ecoji -e
🧜🎺🥧🗽🍻🏺💨🥿🍚📇🌱👞👻🌁🥉🗾
$ echo 🧜🎺🥧🗽🍻🏺💨🥿🍚📇🌱👞👻🌁🥉🗾 | ecoji -d | openssl base64
Qo7e3rIs0pdfySSfYaWNaoO+ZrM=
$ cat encode.go  | openssl dgst -binary -sha1 | openssl base64
Qo7e3rIs0pdfySSfYaWNaoO+ZrM=
```

(If you want to use Ecoji for hashes, consider the dangers inherent in older systems without utf8 emoji support, different fonts, and similar emojis.)

### A URL Shortener

Four base1024 emojis can represent 1 trillion unique IDs.  In the example below `af82dd48f7` represents a 5 byte id for a URL in a key value store like [Accumulo](https://accumulo.apache.org).  When someone enters the URL, the 5 byte id could be used to obtain the actual URL from the database and then redirect.

```
$ printf "https://ecoji.io/%s\n" $(echo af82dd48f7 | xxd -r -p | ecoji -e)
https://ecoji.io/😉🥸🪴🏄
```

## Other Implementations

Libraries implementing the Ecoji encoding standard. Submit PR to add a library to the table.

| Language | Version | Comments |
|----------| ------- | -------- |
| [D](https://github.com/ohdatboi/ecoji-d) | V1 | Implementation of Ecoji written in the D programming language. |
| Go | V2 | This repository offers a Go library package with three functions [ecoji.Encode()](encode.go) [ecoji.Encode2()](encode.go) and [ecoji.Decode()](decode.go). |
| [Java](https://github.com/netvl/ecoji-java) | V1 | Implementation of Ecoji written in Java, usable in any JVM language. |
| [.NET](https://github.com/abock/dotnet-ecoji) | V1 | Implementation of Ecoji written in C# targeting .NET Standard 2.0: [`dotnet add package Ecoji`](https://www.nuget.org/packages/Ecoji). |
| [PHP](https://github.com/Rayne/ecoji-php) | V1 | PHP 7.x implementation of Ecoji. Available as [`rayne/ecoji` on Packagist](https://packagist.org/packages/rayne/ecoji). |
| [Python](https://github.com/mecforlove/ecoji-py) | V1 | Implementation of Ecoji written in the Python3 programming language. |
| [Rust](https://github.com/netvl/ecoji.rs) | V1 | Implementation of Ecoji written in the Rust programming language. |
| [Swift](https://github.com/Robindiddams/ecoji-swift) | V1 | Implementation of Ecoji written in the Swift programming language. |


## Building

If you are new to Go I would recommend this [video] and the [tour].

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
