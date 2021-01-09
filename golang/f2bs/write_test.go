package main

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	// customer
	infn := "/tmp/.."
	oufn := "/tmp/.."
	in, err := os.Open(infn)
	if err != nil {
		t.Fatal(err)
	}
	out, err := os.Create(oufn)
	if err != nil {
		t.Fatal(err)
	}
	err = Write(out, in, false, "build-1", "", "var-1")
	if err != nil {
		t.Fatal(err)
	}
}
