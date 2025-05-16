package errs

import (
	"errors"
	"fmt"
)

type Err struct {
	Code int
	Msg  string
	Data any
	e    error
}

func (e *Err) Error() string {
	return e.Msg
}

func Business(err any) error {
	switch e := err.(type) {
	case *Err:
		return e
	case string:
		return &Err{Code: 600, Msg: e}
	case error:
		return &Err{Code: 600, Msg: e.Error(), e: e}
	default:
		return &Err{Code: 500, Msg: "internal error"}
	}
}

func Internal(err error) error {
	var e *Err
	switch {
	case errors.As(err, &e):
		return e
	default:
		return &Err{Code: 500, Msg: "internal error", e: e}
	}
}

func Param(s any) *Err {
	switch s := s.(type) {
	case string:
		return &Err{Code: 400, Msg: s}
	case error:
		return &Err{Code: 400, Msg: s.Error(), e: s}
	}
	return &Err{Code: 400, Msg: fmt.Sprintf("invalid parameter: %v", s)}
}
