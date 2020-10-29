package defs

type UserCredential struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
}

type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type VideoInfo struct {
	Id           string
	AuthorId     int
	Name         string
	DisplayCtime string
}

type Comment struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

type SimpleSeesion struct {
	Username string
	TTL      int64
}
