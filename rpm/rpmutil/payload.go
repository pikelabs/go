package rpmutil

import (
	"code.pikelabs.net/go/archive/cpio"
)

type PayloadReader interface {
	Next() error
	Read(d []byte) (int, error)
}

type payloadReader struct {
	r cpio.Reader
}
