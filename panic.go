package errors

type internalPanic error

func PanicOnError(err error) {
	if err != nil {
		panic(internalPanic(err))
	}
}

func RecoverError(pErr interface{}) {
	var p = pErr.(*error)
	var r = recover()
	switch r := r.(type) {
	case internalPanic:
		*p =  error(r)
	case nil:
	default:
		panic(r)
	}
}