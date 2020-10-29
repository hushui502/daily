package user

import "awesomeProject2/jianyu/unittest/person"

type User struct {
	Person person.Male
}

func NewUser(p person.Male) *User {
	return &User{Person: p}
}

func (u *User) GetUserInfo(id int64) error {
	return u.Person.Get(id)
}

//mockgen -source=jianyu/unittest/person/male.go -destination=jianyu/unittest/mock/male_mock.go
