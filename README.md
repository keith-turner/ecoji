# ecoji.io source code

Source code for ecoji.io.  The web page calls the ecoji go project using web assembly.  The following builds the web assembly file. 

```
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

