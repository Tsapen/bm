package bm

import (
	"fmt"
)

// Error implements error interface.
type Error string

func (err Error) Error() string {
	return string(err)
}

// ErrPair contains deferred and returned error.
type ErrPair struct {
	Def error
	Ret error
}

// Error returns concatenated error.
func (errPair ErrPair) Error() string {
	return fmt.Sprintf("returned: %s; deferred: %s", errPair.Def, errPair.Ret)
}

// HandleErrPair contains deferred and returned errors.
func HandleErrPair(def, ret error) error {
	if ret == nil {
		return def
	}

	if def == nil {
		return ret
	}

	return ErrPair{
		Def: def,
		Ret: ret,
	}
}
