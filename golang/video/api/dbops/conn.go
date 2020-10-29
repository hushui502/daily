package dbops

import "database/sql"

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("", "")
	if err != nil {
		panic(err.Error())
	}
}
