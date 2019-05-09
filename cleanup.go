package errors

import "log"

// Finalizer collects stacks of functions, which are executed at the end of user defined scopes.
type Finalizer struct {
	parent  *Finalizer
	onScope callStack
	onFinal callStack
}

// NewFinalizer creates a new empty finalizer with no parent.
func NewFinalizer() *Finalizer {
	return new(Finalizer)
}

// ChildScope returns a new finalizer, with the receiver as parent.
func (c *Finalizer) ChildScope() *Finalizer {
	return &Finalizer{parent: c}
}

// EndScope runs all functions added using OnScopeEnd in reverse order.
// If any of the functions return an error or the error referenced by the
// err argument is non-nil, functions added using OnFinalEnd are also called
// in reverse order.
// Any errors encountered are merged with the value referenced by the provided
// error pointer. If the provided error is nil, all encountered errors are logged.
// The empty finalizer may be reused after calling EndScope.
func (c *Finalizer) EndScope(err *error) {
	var tmpErr error
	for c.onScope != nil {
		tmpErr = WithAnother(tmpErr, c.onScope.pop())
	}
	if tmpErr == nil && (err == nil || *err == nil) && c.parent != nil {
		c.parent.OnFinalEnd(c.onFinal...)
	} else {
		for c.onFinal != nil {
			tmpErr = WithAnother(tmpErr, c.onFinal.pop())
		}
	}
	if err != nil {
		*err = WithAftermath(*err, tmpErr)
	} else if tmpErr != nil {
		log.Println("finalizer error:", tmpErr)
	}
}

type callStack []func() error

func (c *callStack) pop() error {
	var l = len(*c)
	var f = (*c)[l-1]
	*c = (*c)[:l-1]
	if len(*c) == 0 {
		*c = nil
	}
	return f()
}

// OnScopeEnd adds a function to be run at the end of this scope (see EndScope).
func (c *Finalizer) OnScopeEnd(f ...func() error) {
	c.onScope = append(c.onScope, f...)
}

// OnFinalEnd adds a function to be run at the end of this scope (see EndScope),
// if an error was encountered. Otherwise the added functions are added to the parent
// finalizer using its OnFinalEnd method.
func (c *Finalizer) OnFinalEnd(f ...func() error) {
	c.onFinal = append(c.onFinal, f...)
}
