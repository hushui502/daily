package v7

import (
	"bufio"
	"io"
)

func CountLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	lines := 0

	for sc.Scan() {
		lines++
	}

	return lines, sc.Err()
}
