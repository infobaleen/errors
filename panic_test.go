package errors

import (
	"github.com/matryer/is"
	"testing"
)

func TestRecoverError(t *testing.T) {
	var is = is.New(t)
	var test = func(shouldFail bool, shouldPanic bool) (err error) {
		defer RecoverError(&err)
		if shouldFail {
			PanicOnError(Fmt("error"))
		}
		if shouldPanic {
			panic("panic")
		}
		return
	}

	is.NoErr(test(false, false))
	is.True(test(true, false) != nil)
	is.True(func() (panicked bool) {
		defer func(){
			panicked = recover() != nil
		}()
		test(false, true)
		return
	}())
}