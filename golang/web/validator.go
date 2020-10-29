package main

//import (
//	"fmt"
//	"gopkg.in/go-playground/validator.v9"
//	"log"
//	"reflect"
//)
//type RegisterReq struct {
//	Username       string `validate:"gt=0"`
//	PasswordNew    string `validate:"gt=0"`
//	PasswordRepeat string `validate:"eqfield=PasswordNew"`
//	Email          string `validate:"email"`
//}
//
//func main() {
//	var s = 23
//	t := reflect.TypeOf(s)
//	v := reflect.ValueOf(s)
//	fmt.Println(t, v)
//
//	//validate(req)
//}
//
//func validate(req RegisterReq) error {
//	v := validator.New()
//	err := v.Struct(req)
//	if err != nil {
//		println(err.Error())
//		log.Fatal(err)
//		return err
//	}
//	return nil
//}
//
//var req = RegisterReq{
//	Username:       "hufan",
//	PasswordNew:    "hufan",
//	PasswordRepeat: "fsd",
//	Email:          "fjlj",
//}