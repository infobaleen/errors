package errors

import "log"

type Finalizer struct {
	parent  *Finalizer
	onScope callStack
	onFinal callStack
}

func NewFinalizer() *Finalizer {
	return new(Finalizer)
}

func (c *Finalizer) ChildScope() *Finalizer {
	return &Finalizer{parent: c}
}

func (c *Finalizer) EndScope(err *error) {
	var tmpErr error
	for c.onScope != nil {
		tmpErr = WithAnother(tmpErr, c.onScope.pop())
	}
	if tmpErr == nil && (err == nil || *err == nil) && c.parent != nil {
		c.parent.onFinal = append(c.parent.onFinal, c.onFinal...)
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

func (c *Finalizer) OnScopeEnd(f ...func() error) {
	c.onScope = append(c.onScope, f...)
}

func (c *Finalizer) OnFinalEnd(f ...func() error) {
	c.onFinal = append(c.onFinal, f...)
}
