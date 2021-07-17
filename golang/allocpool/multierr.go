package allocpool

import "strings"

type MultiError struct {
	errs []error
}

func NewMultiError() *MultiError {
	return &MultiError{}
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.errs = append(m.errs, err)
	}
}

func (m *MultiError) Return() error {
	if len(m.errs) == 0 {
		return nil
	}

	return m
}

func (m *MultiError) Error() string {
	res := strings.Builder{}
	for _, e := range m.errs {
		res.WriteString(e.Error() + "\n")
	}

	return strings.TrimRight(res.String(), "\n")
}
