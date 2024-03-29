package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type errorWithTrace struct {
	pc    [8]uintptr
	cause error
}

func (err errorWithTrace) Error() string {
	var frames = runtime.CallersFrames(err.pc[:])
	var trace strings.Builder
	trace.WriteString("[")
	for {
		var frame, more = frames.Next()
		fmt.Fprint(&trace, path.Base(frame.Function), ":", frame.Line)
		if !more {
			fmt.Fprint(&trace, "] ", err.cause)
			return trace.String()
		}
		trace.WriteString(",")
	}
}

func (err errorWithTrace) Cause() error {
	return err.cause
}

func trace(cause error, skip int) error {
	if cause == nil {
		return nil
	}
	var err errorWithTrace
	err.cause = cause
	runtime.Callers(3+skip, err.pc[:])
	return err
}

func WithTrace(err error) error {
	return trace(err, 0)
}

func WithTraceSkip(err error, n int) error {
	return trace(err, n)
}

func Fmt(format string, args ...interface{}) error {
	return trace(fmt.Errorf(format, args...), 0)
}

type errorWithCause struct {
	message string
	cause   error
}

func Wrap(cause error, messageFormat string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	return trace(errorWithCause{fmt.Sprintf(messageFormat, args...), cause}, 0)
}

func (err errorWithCause) Error() string {
	return fmt.Sprintf("%s: %s", err.message, err.cause)
}

func (err errorWithCause) Cause() error {
	return err.cause
}

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

type errorList []error

func WithAnother(err error, another ...error) error {
	for _, a := range another {
		if a == nil {
			continue
		}
		if err == nil {
			err = a
		} else if list, ok := err.(errorList); ok {
			err = append(list, a)
		} else {
			err = errorList{err, a}
		}
	}
	return err
}

func (err errorList) Error() string {
	var errs []string
	for _, err := range err {
		if s := err.Error(); s != "" {
			errs = append(errs, s)
		}
	}
	return strings.Join(errs, ", ")
}

type errorWithAftermath struct {
	original  error
	aftermath error
}

func WithAftermath(err error, followUp ...error) error {
	if err == nil {
		err = WithAnother(nil, followUp...)
	} else if withAftermath, ok := err.(errorWithAftermath); ok {
		withAftermath.aftermath = WithAnother(withAftermath.aftermath, followUp...)
		err = withAftermath
	} else {
		err = errorWithAftermath{err, WithAnother(nil, followUp...)}
	}
	return err
}

func (err errorWithAftermath) Error() string {
	return fmt.Sprintf("{original: %s; aftermath: %s}", err.original, err.aftermath)
}

func (err errorWithAftermath) Cause() error {
	return err.original
}
