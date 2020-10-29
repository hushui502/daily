package main

import (
	"io/ioutil"
)

func main() {

}

func write() error {
	if err := ioutil.WriteFile("", []byte("fd"), 0644); err != nil {
		return err
	}
	return nil
}

//func read() error {
//	data, err := ioutil.ReadFile("")
//	if err != nil {
//		return err
//	}
//	if err := cfg.Decode(data); err != nil {
//		return err
//	}
//	return nil
//}
