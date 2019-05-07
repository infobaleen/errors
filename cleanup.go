package errors

import (
	"log"
)

type Finalizer struct {
	error  []func() error
	always []func() error
}

func NewFinalizer() *Finalizer {
	return new(Finalizer)
}

func (c *Finalizer) Finalize(err *error) {
	if err != nil && *err != nil {
		callStack(&c.error)
	}
	callStack(&c.always)
}

func callStack(stack *[]func() error) {
	for l := len(*stack); l > 0; l-- {
		var f = (*stack)[l-1]
		*stack = (*stack)[:l-1]
		var err = f()
		if err != nil {
			log.Println("cleanup error: ", err.Error())
		}
	}
}

func (c *Finalizer) ErrFn(f func() error) {
	c.error = append(c.error, f)
}

func (c *Finalizer) FinalFn(f func() error) {
	c.error = append(c.always, f)
}
