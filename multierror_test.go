package multierror

import (
	"fmt"
	"testing"
)

func Works() error { return nil }
func Fails() error { return fmt.Errorf("Failed") }

func FailMulti(msg string, n int) error {
	var e MultiError
	for i := 0; i < n; i++ {
		e.Pushf("%s%d", msg, i+1)
	}
	return e.Error()
}

func ReturnsEmptyMultiError() error {
	var e MultiError
	e.Push(Works())
	e.Push(Works())
	return e.Error()
}

func TestEmptyMultiError(t *testing.T) {
	var e MultiError
	e.Push(Works())
	e.Push(Works())
	if e != nil {
		t.Fatal("Empty MultiError is not nil")
	}

	if ReturnsEmptyMultiError() != nil {
		t.Fatal("Empty MultiError return value is not nil")
	}
}

func ReturnsNonEmptyMultiError() error {
	var e MultiError
	e.Push(Fails())
	e.Push(Fails())
	return e.Error()
}

func TestNonEmptyMultiError(t *testing.T) {
	var e MultiError
	e.Push(Fails())
	e.Push(Fails())
	if e == nil {
		t.Fatal("Non-empty MultiError is nil")
	}

	if ReturnsNonEmptyMultiError() == nil {
		t.Fatal("Non-empty MultiError return value is nil")
	}
}

func TestPushingMultiError(t *testing.T) {
	var e MultiError
	e.Push(FailMulti("Fail", 2))
	if fmt.Sprint(e) != `2 errors: ["Fail1", "Fail2"]` {
		t.Fatal("Incorrect error string: ", e)
	}

	e.Push(FailMulti("X", 3))
	if fmt.Sprint(e) != `5 errors: ["Fail1", "Fail2", "X1", "X2", "X3"]` {
		t.Fatal("Incorrect error string: ", e)
	}
}

func TestMultiErrorStringification(t *testing.T) {
	var e MultiError
	if fmt.Sprint(e) != `nil` {
		t.Fatal("Incorrect error string: ", e)
	}

	e.Push(Fails())
	if fmt.Sprint(e) != `Failed` {
		t.Fatal("Incorrect error string: ", e)
	}

	e.Push(Fails())
	if fmt.Sprint(e) != `2 errors: ["Failed", "Failed"]` {
		t.Fatal("Incorrect error string: ", e)
	}

	if ReturnsNonEmptyMultiError().Error() != `2 errors: ["Failed", "Failed"]` {
		t.Fatal("Incorrect error string: ", ReturnsNonEmptyMultiError())
	}
}
