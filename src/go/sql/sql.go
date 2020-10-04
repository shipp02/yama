package sql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
)

func Get(stmt *sqlx.Stmt, dest interface{}, args ...interface{}) (err error) {
	err = stmt.Get(dest, args)
	if err != nil {
		return
	}
	return nil
}

func Exec(stmt *sqlx.Stmt, args ...interface{}) (err error, res sql.Result) {
	res, err = stmt.Exec(args)
	if err != nil {
		return err, nil
	}
	return
}

// @param dest must be pointer to slice
func Select(stmt *sqlx.Stmt, dest interface{}, args ...interface{}) (err error) {
	err = stmt.Select(dest, args)
	if err != nil {
		return
	}
	return nil
}

func CreateStmts(db *sqlx.DB, sqlStmts []string) (stmts map[string]*sqlx.Stmt) {
	stmts = make(map[string]*sqlx.Stmt, len(sqlStmts))
	for _, sqlStmt := range sqlStmts {
		var err error
		stmts[sqlStmt], err = db.Preparex(sqlStmt)
		if err != nil {
			log.Fatalf("Creation of stmt %s failed due to %s",
				sqlStmt, err.Error())
		}
	}
	return
}
