package utils

import (
	"fmt"
	"runtime/debug"
)

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err, "\n", string(debug.Stack()))
			}
		}()

		f()
	}()
}
