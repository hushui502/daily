package main

import (
	"encoding/json"
	"os"
	"time"
)

type MyUser struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
}

func main() {
	_ = json.NewEncoder(os.Stdout).Encode(
		&MyUser{ID: 1, Name: "ken", LastSeen: time.Now()})
}

//func (u *MyUser) MarshalJSON() ([]byte, error) {
//	return json.Marshal(&struct {
//		ID       int64  `json:"id"`
//		Name     string `json:"name"`
//		LastSeen int64  `json:"lastSeen"`
//	}{
//		ID:       u.ID,
//		Name:     u.Name,
//		LastSeen: u.LastSeen.Unix(),
//	})
//}

func (u *MyUser) MarshalJSON() ([]byte, error) {
	type Alias MyUser
	return json.Marshal(&struct {
		lastSeen int64  `json:"lastSeen"`
		Name     string `json:"name,omitempty"`
		*Alias
	}{
		lastSeen: u.LastSeen.Unix(),
		Alias:    (*Alias)(u),
	})
}
