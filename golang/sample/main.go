package main

import "sync"

type User struct {
	Name string
}

func main() {
	result := make(chan *User)
	var waitGroup sync.WaitGroup
	var users []*User

	for i := 0; i < 100; i++ {
		users = append(users, &User{Name: "test"})
	}

	waitGroup.Add(len(users))
	for _, feed := range users {
		go func(*User) {
			AddUser(feed, result)
			waitGroup.Done()
		}(feed)
	}

	go func() {
		waitGroup.Wait()
		close(result)
	}()

	for res := range result {
		println(res.Name)
	}

}

func AddUser(feed *User, res chan<- *User) {
	res <- feed
}
