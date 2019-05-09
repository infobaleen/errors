package errors

import "log"

type Finalizer struct {
	error  []func() error
	always []func() error
}

func NewFinalizer() *Finalizer {
	return new(Finalizer)
}

func (c *Finalizer) Finalize(err *error) {
	var internalError = callStack(&c.always)
	if internalError != nil {
		internalError = WithAftermath(internalError, callStack(&c.error))
	} else if err != nil && *err != nil {
		internalError = WithAnother(internalError, callStack(&c.error))
	}
	if err != nil {
		*err = WithAftermath(*err, internalError)
	} else if internalError != nil {
		log.Println("finalizer error:", internalError)
	}
}

func callStack(stack *[]func() error) error {
	var err error
	for l := len(*stack); l > 0; l-- {
		err = WithAnother(err, (*stack)[l-1]())
	}
	*stack = nil
	return err
}

func (c *Finalizer) ErrFn(f ...func() error) {
	c.error = append(c.error, f...)
}

func (c *Finalizer) FinalFn(f ...func() error) {
	c.error = append(c.always, f...)
}
