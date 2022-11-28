package main

import "strings"
import "bytes"
import "github.com/keith-turner/ecoji/v2"
import "syscall/js"

func encode(this js.Value, p []js.Value) interface{} {
	reader := strings.NewReader(p[0].String())
	writer := new(bytes.Buffer)
	ecoji.EncodeV2(reader, writer, 0)
	return js.ValueOf(writer.String())
}

func decode(this js.Value, p []js.Value) interface{} {
	reader := strings.NewReader(p[0].String())
	writer := new(bytes.Buffer)
	err := ecoji.Decode(reader, writer)
	// TODO ensure output is UTF-8
	if err != nil {
		return js.ValueOf("🤨 It seems that your input was not Ecoji™ encoded")
	}
	return js.ValueOf(writer.String())
}

func main() {
	done := make(chan struct{})
	js.Global().Set("ecojiEncode", js.FuncOf(encode))
	js.Global().Set("ecojiDecode", js.FuncOf(decode))
	<-done
}
