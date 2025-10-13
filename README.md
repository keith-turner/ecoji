# [ecoji.io](https://ecoji.io) source code

Source code for ecoji.io. All of the web page source is in [index.html](index.html) and [ecoji.js](ecoji.js).  The `ecoji.wasm` file was created by building [ecoji-zig](https://codeberg.org/keith-turner/ecoji-zig).

For local testing, build and run [ths.go](ths.go) for a simple local web server.

```bash
go build -o ths ths.go
./ths &
# open localhost:8080 in a local web browser
# kill process when done
kill %1
```
