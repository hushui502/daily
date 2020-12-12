package proxy

import (
	"log"
	"time"
)

type IUser interface {
	Login(username, password string) error
}

type User struct {
	username string
	password string
}

func (u User) Login(username, password string) error {
	return nil
}

type UserProxy struct {
	user *User
}

func NewUserProxy(user *User) *UserProxy {
	return &UserProxy{
		user: user,
	}
}


// 似乎违反SRP，但是放大一点讲这就是整个登录功能，也可以解释，这样避免了大量的时序耦合的重复代码
func (p *UserProxy) Login(username, password string) error {
	// before 登录前的逻辑
	start := time.Now()

	if err := p.user.Login(username, password); err != nil {
		return err
	}

	// 登录后的逻辑
	log.Printf("user login cost time: %s", time.Now().Sub(start))

	return nil

}