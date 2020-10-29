package sub2

import (
	"github.com/pkg/errors"
	"io/ioutil"
)

func Diff(foo int, bar int) error {
	return errors.New("sub diff error")
}

func IoDiff(foo int, bar int) ([]byte, error) {
	b, err := ioutil.ReadFile("filename")
	return b, err
}

