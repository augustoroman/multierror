package multierror

import (
	"fmt"
	"testing"
)

func Works() error { return nil }
func Fails() error { return fmt.Errorf("Failed") }

func FailAccumulator(msg string, n int) error {
	var e Accumulator
	for i := 0; i < n; i++ {
		e.Pushf("%s%d", msg, i+1)
	}
	return e.Error()
}

func ReturnsEmptyAccumulator() error {
	var e Accumulator
	e.Push(Works())
	e.Push(Works())
	return e.Error()
}

func TestEmptyAccumulator(t *testing.T) {
	var e Accumulator
	e.Push(Works())
	e.Push(Works())
	if e != nil {
		t.Fatal("Empty Accumulator is not nil")
	}

	if ReturnsEmptyAccumulator() != nil {
		t.Fatal("Empty Accumulator return value is not nil")
	}
}

func ReturnsNonEmptyAccumulator() error {
	var e Accumulator
	e.Push(Fails())
	e.Push(Fails())
	return e.Error()
}

func TestNonEmptyAccumulator(t *testing.T) {
	var e Accumulator
	e.Push(Fails())
	e.Push(Fails())
	if e == nil {
		t.Fatal("Non-empty Accumulator is nil")
	}

	if ReturnsNonEmptyAccumulator() == nil {
		t.Fatal("Non-empty Accumulator return value is nil")
	}
}

func TestPushingAccumulator(t *testing.T) {
	var e Accumulator
	e.Push(FailAccumulator("Fail", 2))
	if fmt.Sprint(e) != `2 errors: ["Fail1", "Fail2"]` {
		t.Fatal("Incorrect error string: ", e)
	}

	e.Push(FailAccumulator("X", 3))
	if fmt.Sprint(e) != `5 errors: ["Fail1", "Fail2", "X1", "X2", "X3"]` {
		t.Fatal("Incorrect error string: ", e)
	}
}

func TestAccumulatorStringification(t *testing.T) {
	var e Accumulator
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

	if ReturnsNonEmptyAccumulator().Error() != `2 errors: ["Failed", "Failed"]` {
		t.Fatal("Incorrect error string: ", ReturnsNonEmptyAccumulator())
	}
}
