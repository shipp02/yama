package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingcap/errors"
	"log"
	"testing"
)
import _ "github.com/mattn/go-sqlite3"

const schema = `
CREATE TABLE sql_go_test (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(256),
    username VARCHAR(256),
    password VARCHAR(256)
)
`
const (
	GET    = "SELECT * FROM sql_go_test"
	INSERT = "INSERT INTO sql_go_test (name, username, password) VALUES (?, ?, ?)"
)

type context struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func Setup() *context {
	var c *context
	setup := func() {
		var err error
		c = new(context)
		c.db, err = sqlx.Connect("sqlite3", "file::memory:")
		if err != nil {
			log.Fatalf("Failed to connect to database %s", err)
			return
		}
		_, err = c.db.Exec(schema)
		if err != nil {
			log.Fatalf("Failed to create table %s", err)
			return
		}
	}
	if c == nil {
		setup()
		return c
	}
	return c
}

func TestGet(t *testing.T) {

}

func TestCreateStmts(t *testing.T) {
	stmts := make([]string, 2)
	c := Setup()
	stmts[0] = GET
	stmts[1] = INSERT
	var err error
	c.stmts, err = CreateStmts(c.db, stmts)
	if err != nil {
		for err := errors.Unwrap(err); err != nil; {
			log.Println(err.Error())
		}
		t.Fail()
		return
	}
}

//func TestExec(t *testing.T) {
//	c := Setup()
//	err, _ := Exec(c.stmts[INSERT], "aashay", "sanjay", "pass")
//	if err != nil {
//		t.Fail()
//		return
//	}
//
//}
