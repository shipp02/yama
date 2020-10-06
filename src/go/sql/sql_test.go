package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pingcap/errors"
	"log"
	"sync"
	"testing"
)

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
	free  func()
}

func Setup(t *testing.T) *context {
	var c *context
	var wg sync.WaitGroup
	setup := func() {
		var err error
		c = new(context)
		c.db, err = sqlx.Connect("sqlite3", "file::memory:")
		//        c.db, err = sqlx.Connect("mysql", "root:yousql@tcp(127.0.0.1:3306)/test")
		if err != nil {
			t.Logf("Failed to connect to database %s", err)
			return
		}
		_, err = c.db.Exec(schema)
		if err != nil {
			t.Logf("Failed to create table %s", err)
			return
		}
		c.free = func() {
			t.Log("Freed data")
			wg.Done()
		}
		log.Println("Data Source Created")
		go func() {
			wg.Wait()
			_ = c.db.Close()
		}()
	}
	if c == nil {
		setup()
	}
	wg.Add(1)
	return c
}

func TestGet(t *testing.T) {

}

func TestCreateStmts(t *testing.T) {
	stmts := make([]string, 2)
	c := Setup(t)
	stmts[0] = GET
	stmts[1] = INSERT
	var err error
	c.stmts, err = CreateStmts(c.db, stmts)
	if err != nil {
		for err := errors.Unwrap(err); err != nil; {
			t.Log(err.Error())
		}
		t.Fail()
		return
	}
	c.free()
}

func TestExec(t *testing.T) {
	c := Setup(t)
	//	if len(c.stmts) == 0 {
	//		c = Setup()
	//	}
	stmt, err := c.db.Preparex(INSERT)
	if err != nil {
		t.Log("Statement Creation Failed")
		t.Fail()
		return
	}
	args := []interface{}{"aashay", "sanjay", "pass"}
	err, _ = Exec(stmt, args...)
	if err != nil {
		t.Log("TestExec execution failed", err)
		t.Fail()
		return
	}
	c.free()
}
