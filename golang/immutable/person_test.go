package main

import "testing"

type personMock struct {
	Person
	receivedNewColor string
}

func (m personMock) WithFavoriteColorAt(i int, favoriteColor string) Person {
	m.receivedNewColor = favoriteColor
	return m
}


func TestPerson_WithFavoriteColorAt(t *testing.T) {
	mock := personMock{}
	result := updateFavoriteColors(mock)
}
