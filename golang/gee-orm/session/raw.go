package session

import (
	"awesomeProject2/gee-orm/clause"
	"awesomeProject2/gee-orm/dialect"
	"awesomeProject2/gee-orm/log"
	"awesomeProject2/gee-orm/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db			*sql.DB
	tx 			*sql.Tx
	sql			strings.Builder
	sqlVals 	[]interface{}
	refTable 	*schema.Schema
	dialect 	dialect.Dialect
	clause 		clause.Clause
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:db,
		dialect:dialect,
	}
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return 
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	return s.DB().QueryRow(s.sql.String(), s.sqlVals...)
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVals = append(s.sqlVals, values...)
	return s
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVals = nil
	s.clause = clause.Clause{}
}
