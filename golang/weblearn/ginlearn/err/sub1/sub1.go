package sub1

import (
	"awesomeProject2/ginlearn/err/sub1/sub2"
	"github.com/pkg/errors"
)

func Diff(foo int, bar int) error {
	if foo < 0 {
		return errors.New("diff error")
	}
	if err := sub2.Diff(foo, bar); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func IoDiff(foo, bar int) error {
	_, err := sub2.IoDiff(foo, bar)
	return err
}

