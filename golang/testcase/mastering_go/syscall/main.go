package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	fmt.Println("User id: ", os.Getuid())
	var u *user.User
	u, _ = user.Current()
	fmt.Println("Group ids: ")
	groupIds, _ := u.GroupIds()
	for _, i := range groupIds {
		fmt.Print(i, " ")
	}
	fmt.Println(" ")
}
