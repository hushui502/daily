package models

import (
	"gopkg.in/mgo.v2/bson"
	"honeypot/manager/util"
)

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	UserName string
	Password string
}

func ListUser() (users []User, err error) {
	err = collAdmin.Find(nil).All(&users)
	return users, err
}

func GetUserById(id string) (user User, err error) {
	err = collAdmin.FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func NewUser(username, password string) (err error) {
	encryptPass := util.EncryptPass(password)
	err = collAdmin.Insert(&User{Id: bson.NewObjectId(), UserName: username, Password: encryptPass})
	return err
}

func UpdateUser(id string, username, password string) (err error) {
	user := new(User)
	err = collAdmin.FindId(bson.ObjectIdHex(id)).One(user)
	user.UserName = username
	user.Password = util.EncryptPass(password)
	err = collAdmin.UpdateId(bson.ObjectIdHex(id), user)
	return err
}

func DelUser(id string) (err error) {
	err = collAdmin.RemoveId(bson.ObjectIdHex(id))
	return err
}

func Auth(username, password string) (result bool, err error) {
	encryptPass := util.EncryptPass(password)
	userAuth := User{}
	err = collAdmin.Find(bson.M{"username": username, "password": encryptPass}).One(&userAuth)
	if err == nil && userAuth.UserName == username && userAuth.Password == encryptPass {
		result = true
	}
	return result, err
}
