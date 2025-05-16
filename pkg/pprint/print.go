package pprint

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Print(val interface{}) {
	fmt.Println(toString(val))
}

func Println(val interface{}) {
	fmt.Println(toString(val))
}

func toString(val interface{}) string {
	return string(toBytes(val))
}

func toBytes(val interface{}) []byte {
	buf := bytes.Buffer{}

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.Encode(val)

	return buf.Bytes()
}
