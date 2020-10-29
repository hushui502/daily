package main

import "errors"

type Option struct {
	Key string
}

var mux map[string]func(option *Option) error

func register(key string, f func(option *Option) error) error {
	if mux == nil {
		mux = make(map[string]func(option *Option) error)
	}
	if _, exist := mux[key]; exist {
		return errors.New("handler exist")
	}
	mux[key] = f
	return nil
}

func factory(option *Option) error {
	return mux[option.Key](option)
}
