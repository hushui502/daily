package session

import (
	"awesomeProject2/gee-orm/dialect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)
var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)


func NewSession() *Session {
	TestDB, _ = sql.Open("sqlite3", "../main.db")
	return New(TestDB, TestDial)
}

//func TestMain(m *testing.M) {
//	TestDB, _ := sql.Open("sqlite3", "../main.db")
//	code := m.Run()
//	_ = TestDB.Close()
//	os.Exit(code)
//}

func TestSession_Exec(t *testing.T) {
	TestDB, _ = sql.Open("sqlite3", "../main.db")
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	TestDB, _ = sql.Open("sqlite3", "../main.db")
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}

