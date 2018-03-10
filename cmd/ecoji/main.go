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

Options:
    -d, --decode          decode data
    -w, --wrap=COLS       wrap encoded lines after COLS character (default 76).
                          Use 0 to disable line wrapping
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

	decode := false
	wrap := uint(76)

	flag.BoolVar(&decode, "d", false, "")
	flag.BoolVar(&decode, "decode", false, "")

	flag.UintVar(&wrap, "w", 76, "")
	flag.UintVar(&wrap, "wrap", 76, "")

	flag.Usage = func() {
		fmt.Print(usageMessage)
	}

	flag.Parse()

	args := flag.Args()

	if len(args) > 1 {
		fmt.Print(usageMessage)
		os.Exit(2)
	}

	var in *bufio.Reader

	if len(args) == 0 {
		in = bufio.NewReader(os.Stdin)
	} else {
		f := openFile(args[0])
		in = bufio.NewReader(f)
	}

	stdout := bufio.NewWriter(os.Stdout)

	if !decode {
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
