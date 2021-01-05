package main

import "os"

func main() {
	ParallelWrite([]byte("xxx"))
}

func ParallelWrite(data []byte) chan error {
	res := make(chan error, 2)
	f1, err := os.Create("/tmp/file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			// if _, err = f1.Write(data) ==>
			// This err is shared with the main goroutine,
			// so the write races with the write below.
			_, err := f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("/tmp/file2") // The second conflicting write to err.
	if err != nil {
		res <- err
	} else {
		go func() {
			// if _, err = f2.Write(data) ==> ...
			// so we should ues _, err := f2.Write(data) instead of ...
			_, err := f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	return res
}
