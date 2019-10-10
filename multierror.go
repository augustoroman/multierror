// Package multierror defines an error Accumulator to contain multiple errors.
package multierror

import (
	"bytes"
	"fmt"
)

// Accumulator is an error accumulator.
//
// Usage:
//
//     var errors multierror.Accumulator
//     errors.Push(returnsErrOrNil())
//     errors.Push(returnsErrOrNil())
//     errors.Push(returnsErrOrNil())
//     return errors.Error()
type Accumulator []error

// Push adds an error to the Accumulator.  If err is nil, then Accumulator is
// not affected.
func (m *Accumulator) Push(err error) {
	if err == nil {
		return
	}
	// Check for a Accumulator
	if e, ok := err.(_error); ok {
		*m = append(*m, e...)
		return
	}

	*m = append(*m, err)
}

// Pushf adds a formatted error string to the Accumulator.  It is a shortcut
// for Push(fmt.Error(...)).
func (m *Accumulator) Pushf(fmtstr string, args ...interface{}) {
	*m = append(*m, fmt.Errorf(fmtstr, args...))
}

// Error returns the accumulated errors.  If no errors have been pushed onto the
// accumulator, then nil will be returned.
func (m *Accumulator) Error() error {
	if len(*m) == 0 {
		return nil
	} else if len(*m) == 1 {
		return (*m)[0]
	}
	return _error(*m)
}

// String prints the accumulated errors or "nil" if no errors have been pushed.
func (m Accumulator) String() string {
	return _error(m).Error()
}

// This type implements the actual error interface.  This is separate from
// Accumulator so to avoid accidentally returning an interface to a nil pointer.
type _error []error

func (m _error) Error() string {
	if len(m) == 0 {
		return "nil"
	}
	if len(m) == 1 {
		return m[0].Error()
	}

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d errors: [", len(m))
	for i, err := range m {
		if i != 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, `"%v"`, err)
	}
	fmt.Fprint(buf, "]")
	return buf.String()
}
