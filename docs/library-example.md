# Using Ecoji as a library

This is an example that shows how to use Ecoji as a library.  Start off with a new directory.

```bash
mkdir /tmp/ecoji-ip
cd /tmp/ecoji-ip
```

Then create file named `ecoji-ip.go` with the following contents.  This example program can encode and decode IPv4 addresses using Ecoji V2.

```go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/keith-turner/ecoji/v2"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <ip addr OR ecoji encoded ip>\n", os.Args[0])
		os.Exit(1)
	}

	// attempt to pare an IP address
	ip := net.ParseIP(os.Args[1])

	// this example does not handle IPv6

	if ip != nil {
		outbuf := bufio.NewWriter(os.Stdout)

		//encode the 32 bit or 4 byte representation of the IP address as Ecoji

		if err := ecoji.EncodeV2(bytes.NewReader(ip.To4()), outbuf, 4); err != nil {
			log.Fatal(err)
		}

		outbuf.Flush()
	} else {
		buffer := bytes.NewBuffer(nil)
		if err := ecoji.Decode(strings.NewReader(os.Args[1]), buffer); err != nil {
			log.Fatal("Input is not Ecoji or an IP address.")
		}

		// ensure the encoded data is exactly 32 bits or 4 bytes
		if len(buffer.Bytes()) != 4 {
			log.Fatal("Encoded data does not appear to be an IP address")
		}

		ip = buffer.Bytes()
		fmt.Println(ip.String())
	}
}
```

Then run the following commands to build an executable.

```
go mod init local/ecoji-ip
go mod tidy
go build ecoji-ip.go
```

Now you should be able to encode and decode.  Below is an example.

```
$ ./ecoji-ip 140.82.121.4
🧳🏸🔔🥷
$ ./ecoji-ip 104.244.42.129
👴🚾🪹🛼
$ ./ecoji-ip 🧳🏸🔔🥷
140.82.121.4
$ ./ecoji-ip 👴🚾🪹🛼
104.244.42.129
```
