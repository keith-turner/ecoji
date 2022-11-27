package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/keith-turner/ecoji"
	"log"
	"os"
)

var usageMessage = `usage: ecoji [OPTIONS]... [FILE]

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

🏣🔉🦐🩻🍈🚞🤹🥷
`

var versionMessage = `Ecoji version 2.0.0
  Copyright   : (C) 2021 Keith Turner
  License     : Apache 2.0
  Source code : https://github.com/keith-turner/ecoji
`

func openFile(name string) *os.File {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		//TODO use log.fatal ??
		fmt.Printf("ERROR : %s \n", err.Error())
		os.Exit(2)
	}

	stat, err2 := f.Stat()

	if err2 != nil {
		//TODO use log.fatal ??
		fmt.Printf("ERROR : %s \n", err.Error())
		os.Exit(2)
	}

	if stat.IsDir() {
		fmt.Printf("ERROR : %s is a directory\n", name)
		os.Exit(2)
	}

	return f
}

func main() {

	encode := false
	decode := false
	help := false
	version := false
	wrap := uint(76)

	flag.BoolVar(&encode, "e", false, "")
	flag.BoolVar(&encode, "encode", false, "")

	flag.BoolVar(&decode, "d", false, "")
	flag.BoolVar(&decode, "decode", false, "")

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")

	flag.UintVar(&wrap, "w", 76, "")
	flag.UintVar(&wrap, "wrap", 76, "")

	flag.Usage = func() {
		fmt.Print(usageMessage)
	}

	flag.Parse()

	args := flag.Args()

	if len(args) > 1 || (encode && decode) {
		fmt.Print(usageMessage)
		os.Exit(2)
	}

	if help {
		fmt.Print(usageMessage)
		return
	}

	if version {
		fmt.Print(versionMessage)
		return
	}

	var in *bufio.Reader

	if len(args) == 0 {
		in = bufio.NewReader(os.Stdin)
	} else {
		f := openFile(args[0])
		in = bufio.NewReader(f)
	}

	stdout := bufio.NewWriter(os.Stdout)

	if encode {
		if err := ecoji.EncodeV2(in, stdout, wrap); err != nil {
			log.Fatal(err)
		}
	} else if !decode {
		if err := ecoji.Encode(in, stdout, wrap); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := ecoji.Decode(in, stdout); err != nil {
			log.Fatal(err)
		}
	}

	stdout.Flush()
}
