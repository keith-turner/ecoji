# Using Ecoji as a library

This is an example that shows how to use Ecoji as a library.  Start off with a new directory.

```bash
mkdir /tmp/ecoji-example
cd /tmp/ecoji-example
```

Then create a file called `go.mod` in that directory with the following contents.

```
go 1.18

module local/ecoji-example

require github.com/keith-turner/ecoji/v2 v2.0.0

```

Then create file named `ecoji-example.go` with the following contents.  This example program will read a file, encode with Ecoji V2, and write that to a new file.

```go
package main

import (
	"bufio"
	"fmt"
	"github.com/keith-turner/ecoji/v2"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("Usage : %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	input := os.Args[1]
	output := os.Args[2]

	infile, err1 := os.OpenFile(input, os.O_RDONLY, 0)
	if err1 != nil {
		log.Fatal(err1)
	}

	outfile, err2 := os.Create(output)
	if err2 != nil {
		log.Fatal(err2)
	}

	outbuf := bufio.NewWriter(outfile)

	//encode data using Ecoji V2 with a 72 emoji wrap
	if err := ecoji.EncodeV2(bufio.NewReader(infile), outbuf, 72); err != nil {
		log.Fatal(err)
	}

	infile.Close()
	outbuf.Flush()
	outfile.Close()
}
```

Then run the following commands to build the executable.

```
go mod download github.com/keith-turner/ecoji/v2
go build ecoji-example.go
```

Now you should be able to use the executable to encode a file.  Below is an example if running it.

```
/tmp/ecoji-example$ ./ecoji-example go.mod test.ecoji
/tmp/ecoji-example$ cat test.ecoji 
👮😽♏🧧🎌🤭🪜🕋💎🧵🐬🌭🍉😑🦎🛸💁🙁🐩🤜👺🕺🛫👉👖😢⛲🌭🛞🍯🍡👂💦🪳🍡🏮👮🪳🏨🌽👙😱🎨🤚🎅😁🐫👅👱😢🏫👃💊🔪🍓🤒👞🙁🖖🐀💩🚞⛵🔅🎀🙎🤹♍🛞☕

```
