package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// CREATE TABLE posts(id int NOT NULL AUTO_INCREMENT,owner int NOT NULL,text TEXT,PRIMARY KEY(id), FOREIGN KEY (owner) REFERENCES users (id) );
const PostSchema = `
CREATE TABLE posts(
	id int NOT NULL AUTO_INCREMENT,
	owner int NOT NULL,
	text TEXT,
	PRIMARY KEY(id)
)
`
type Post struct {
	id int64
	ownerID int64
	text string
}

type queryPost struct {
	id sql.NullInt64
	OwnerID sql.NullInt64
	text sql.NullString
}

func (qp *queryPost) GetInterface(l int) (s []interface{}) {
	s = make([]interface{}, l)
	s[0] = &qp.id
	s[1] = &qp.OwnerID
	s[2] = &qp.text
	return
}

func (qp *queryPost) ToPost()(p *Post) {
	// fmt.Println(*qp)
	p = new(Post)
	p.id = qp.id.Int64
	p.ownerID = qp.OwnerID.Int64
	p.text = qp.text.String
	// fmt.Println(p)
	return
}

func GetPost(db *sqlx.DB, p *Post)(*Post, error){
	var err error
	var query = `
		SELECT * 
		FROM posts
	`

	var idQ = "WHERE id=$(ID)"
	var oidQ = "WHERE owner=$(OID)"

	if p.id == 0 &&  p.ownerID == 0 {
		err = errors.New("Insufficient data")
	}

	var where string

	if p.id != 0 && where == ""{
		where = strings.Replace(idQ, "$(ID)", strconv.FormatInt(p.id, 10), 1)
	}
	if p.ownerID != 0 && where == ""{
		where = strings.Replace(oidQ, "$(OID)", strconv.FormatInt(p.ownerID, 10), 1)
	}
	query += where
	resp, err := db.Query(query)
	l,err := resp.Columns()

	var qp  =  new(queryPost)
	is := qp.GetInterface(len(l))

	for resp.Next() {
		if err:=resp.Scan(is...); err !=nil {
			log.Fatal(err)
			err = errors.New(err.Error())
		}
		// fmt.Println(is...)
	}

	p = qp.ToPost()
	// fmt.Printf("%p\n", p)
	return p, err
}

func (p *Post) CreatePost (db *sqlx.DB) (error){
	var err error
	qp, _ := GetPost(db, p)
	if qp.id != 0 {
		err = errors.New("Post exists")
	}

	var exec = "INSERT INTO posts (owner, text) VALUES($(OID), \"$(TEXT)\")"

	exec = strings.Replace(exec, "$(OID)", strconv.FormatInt(p.ownerID, 10), 1)
	exec = strings.Replace(exec, "$(TEXT)", p.text, 1)
	fmt.Println(exec)

	db.MustExec(exec)
	return err
}

func Connect() (db *sqlx.DB){
	db,err:= sqlx.Connect("mysql", "root:yoursql@tcp(localhost:3306)/mysql")
	if err != nil {
        log.Fatalln(err)
    }
	return 
}

func main() {
	db := Connect()
	db.MustExec(PostSchema)
	db.MustExec("INSERT INTO posts (owner, text) VALUES(20, \"NEW POST EH\")")
	db.MustExec("INSERT INTO posts (owner, text) VALUES(21, \"another one EH\")")
	p := new(Post)
	p.ownerID = 10
	p.text = "New method eh"
	p.CreatePost(db)

	p.ownerID = 10
	p,_ = GetPost(db, p)
	// fmt.Printf("%p\n", p)
	fmt.Println(p)
}
