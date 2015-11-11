package main

import (
	"fmt"
	"io"
)

type errWriter struct {
	w   io.Writer
	err error
}

func (e *errWriter) writeByte(p []byte) {
	if e.err != nil {
		return
	}
	fmt.Println(string(p))
	_, e.err = e.w.Write(p)
}

func (e *errWriter) writeString(p string) {
	if e.err != nil {
		return
	}
	// pvec := *(*[]byte)(unsafe.Pointer(&p))
	// _, e.err = e.w.Write(pvec)
	_, e.err = e.w.Write([]byte(p))
}

func (e *errWriter) Err() error {
	return e.err
}
