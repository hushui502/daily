package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func main() {
	var err error
	db1, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/hand1")
	if err != nil {
		panic(err)
	}
	defer db1.Close()

	db2, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/hade2")
	if err != nil {
		panic(err.Error())
	}
	defer db2.Close()


	// 开始前显示
	var score int
	db1.QueryRow("select score from user where id = 1").Scan(&score)
	fmt.Println("user1 score:", score)
	var money float64
	db2.QueryRow("select money from wallet where id = 1").Scan(&money)
	fmt.Println("wallet1 money:", money)

	tx1, err := db1.Begin()
	if err != nil {
		panic(err)
	}
	tx2, err := db2.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			//tx1.Rollback()
			//tx2.Rollback()
		}
	}()

	xid := strconv.FormatInt(time.Now().Unix(), 10)
	if _, err = db1.Exec(fmt.Sprintf("XA START '%s'", xid)); err != nil {
		
	}
	tx1.Exec("update user set score=score+2 where id = 1")


	db1.Exec(fmt.Sprintf("XA END '%s'", xid))

	db1.Exec(fmt.Sprintf("XA PREPARE '%s'", xid))


}



