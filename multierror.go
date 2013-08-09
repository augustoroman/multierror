package multierror

import (
	"bytes"
	"fmt"
)

type MultiError []error

func (m *MultiError) Push(err error) {
	if err == nil {
		return
	}
	// Check for a multierror
	if e, ok := err.(_error); ok {
		*m = append(*m, e...)
		return
	}

	*m = append(*m, err)
}

func (m *MultiError) Pushf(fmtstr string, args ...interface{}) {
	*m = append(*m, fmt.Errorf(fmtstr, args...))
}

func (m *MultiError) Error() error {
	if len(*m) == 0 {
		return nil
	}
	return _error(*m)
}

func (m MultiError) String() string {
	return _error(m).Error()
}

// This type implements the error interface.  This is separate from MultiError
//
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
