package array

import (
	"errors"
	"fmt"
)

type Array struct {
	data []int
	length uint
}

func NewArray(capacity uint) *Array {
	if capacity <= 0 {
		return nil
	}

	// TODO reflect type
	return &Array{
		data:   make([]int, capacity, capacity),
		length: 0,
	}
}

func (arr *Array) Len() uint {
	return arr.length
}

func (arr *Array) isIndexOutOfRange(index uint) bool {
	if index >= uint(cap(arr.data)) {
		return true
	}

	return false
}

func (arr *Array) Find(index uint) (int, error) {
	if arr.isIndexOutOfRange(index) {
		return 0, errors.New("out of index range")
	}

	return arr.data[index], nil
}

func (arr *Array) Insert(index uint, v int) error {
	if arr.Len() == uint(cap(arr.data)) {
		return errors.New("full array")
	}
	if index != arr.length && arr.isIndexOutOfRange(index) {
		return errors.New("out of index array")
	}

	for i := arr.length; i > index; i-- {
		arr.data[i] = arr.data[i-1]
	}

	arr.data[index] = v
	arr.length++

	return nil
}

func (arr *Array) InsertToTail(v int) error {
	return arr.Insert(arr.Len(), v)
}

func (arr *Array) Delete(index uint) (int, error) {
	if arr.isIndexOutOfRange(index) {
		return 0, errors.New("out of index range")
	}

	v := arr.data[index]
	arr.data = append(arr.data[:index], arr.data[index+1:]...)
	arr.length++

	return v, nil
}

func (arr *Array) Print() {
	var format string
	for i := uint(0); i < arr.Len(); i++ {
		format += fmt.Sprintf("|%+v", arr.data[i])
	}

	fmt.Println(format)
}