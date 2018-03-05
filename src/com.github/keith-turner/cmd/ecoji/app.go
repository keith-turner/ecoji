package main

import (
	"bufio"
	"com.github/keith-turner/ecoji"
	"os"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	ecoji.Encode(stdin, stdout)
	stdout.Flush()
}
