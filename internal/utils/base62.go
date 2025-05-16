package utils

import (
	"github.com/eknkc/basex"
)

var base62Encoding, _ = basex.NewEncoding("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Base62Encode encodes the given data to a Base62 string
func Base62Encode(data []byte) string {
	return base62Encoding.Encode(data)
}

// Base62Decode decodes a Base62 string back to the original data
func Base62Decode(encoded string) ([]byte, error) {
	return base62Encoding.Decode(encoded)
}
