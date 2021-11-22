package errors

import (
	"fmt"
	"testing"
)

func TestWithAnother(t *testing.T) {
	var err = WithAnother(fmt.Errorf("err1"), fmt.Errorf("err2"))
	fmt.Println(err.Error())
}
