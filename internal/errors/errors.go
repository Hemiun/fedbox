package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// IncludeBacktrace is a static variable that decides if when creating an error we store the backtrace with it.
var IncludeBacktrace = true

// Err is our custom error type
type Err struct {
	c error
	m string
	t []byte
	l int
	f string
}

func (e Err) Error() string {
	return e.m
}

type wrapper interface {
	Unwrap() error
}

func (e Err) Unwrap() error {
	return e.c
}

func (e *Err) Location() (string, int) {
	return e.f, e.l
}

func (e *Err) StackTrace() []byte {
	return e.t
}

func Annotatef(e error, s string, args ...interface{}) error {
	err := wrap(e, s, args...)
	return &err
}

func Newf(s string, args ...interface{}) error {
	err := wrap(nil, s, args...)
	return &err
}

func wrap(e error, s string, args ...interface{}) Err {
	err := Err{
		c: e,
		m: fmt.Sprintf(s, args...),
	}
	if IncludeBacktrace {
		_, err.f, err.l, _ = runtime.Caller(2)
		err.t = debug.Stack()
	}
	return err
}

func Errorf(s string, args ...interface{}) error {
	err := wrap(nil, s, args...)
	return &err
}
func (e *Err) As(err interface{}) bool {
	switch x := err.(type) {
	case **Err:
		(*x).m = e.m
		(*x).c = e.c
		(*x).t = e.t
		(*x).l = e.l
		(*x).f = e.f
	case *Err:
		x.m = e.m
		x.c = e.c
		x.t = e.t
		x.l = e.l
		x.f = e.f
	default:
		return false
	}
	return true
}
