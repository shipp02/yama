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

const PostSchema = `
CREATE TABLE posts(
	id int NOT NULL AUTO_INCREMENT,
	owner_id int NOT NULL,
	text TEXT,
	PRIMARY KEY(id),
	FOREIGN KEY (owner_id) REFERENCES users (id) 
);
`

type Post struct {
	Id      int64  `db:"id"`
	OwnerID int64  `db:"owner_id"`
	Text    string `db:"text"`
}

type queryPost struct {
	id      sql.NullInt64
	OwnerID sql.NullInt64
	text    sql.NullString
}

func (qp *queryPost) GetInterface(l int) (s []interface{}) {
	s = make([]interface{}, l)
	s[0] = &qp.id
	s[1] = &qp.OwnerID
	s[2] = &qp.text
	return
}

func (qp *queryPost) ToPost() (p *Post) {
	// fmt.Println(*qp)
	p = new(Post)
	p.Id = qp.id.Int64
	p.OwnerID = qp.OwnerID.Int64
	p.Text = qp.text.String
	// fmt.Println(p)
	return
}

func GetPost(db *sqlx.DB, p *Post) (*Post, error) {
	var err error
	var query = `
		SELECT * 
		FROM posts
	`

	var idQ = "WHERE id=$(ID)"
	var oidQ = "WHERE owner_id=$(OID)"

	if p.Id == 0 && p.OwnerID == 0 {
		err = errors.New("Insufficient data")
	}

	var where string

	if p.Id != 0 && where == "" {
		where = strings.Replace(idQ, "$(ID)", strconv.FormatInt(p.Id, 10), 1)
	}
	if p.OwnerID != 0 && where == "" {
		where = strings.Replace(oidQ, "$(OID)", strconv.FormatInt(p.OwnerID, 10), 1)
	}
	query += where
	resp, err := db.Query(query)
	fmt.Println(resp, err)
	l, err := resp.Columns()

	var qp = new(queryPost)
	is := qp.GetInterface(len(l))

	for resp.Next() {
		if err := resp.Scan(is...); err != nil {
			log.Fatal(err)
			err = errors.New(err.Error())
		}
		// fmt.Println(is...)
	}

	p = qp.ToPost()
	// fmt.Printf("%p\n", p)
	return p, err
}

func (u *User) GetPosts(db *sqlx.DB) ([]Post, error){
	var err error
	posts := []Post{}
	query := `
		SELECT * FROM posts
		WHERE owner_id=$(OID)
	`
	query = strings.Replace(query, "$(OID)", strconv.FormatInt(u.id, 10), 1)
	fmt.Println(query)
	erro:=db.Select(&posts, query)
	if erro != nil{
		fmt.Println(erro)
	}
	fmt.Println(posts)
	return posts, err
}

func (p *Post) CreatePost(db *sqlx.DB) error {
	var err error
	qp, _ := GetPost(db, p)
	if qp.Id != 0 {
		err = errors.New("Post exists")
	}

	var exec = "INSERT INTO posts (owner_id, text) VALUES($(OID), \"$(TEXT)\")"

	exec = strings.Replace(exec, "$(OID)", strconv.FormatInt(p.OwnerID, 10), 1)
	exec = strings.Replace(exec, "$(TEXT)", p.Text, 1)
	fmt.Println(exec)

	db.MustExec(exec)
	return err
}

func mainP() {
	db := Connect()
	db.MustExec(PostSchema)
	db.MustExec("INSERT INTO posts (owner, text) VALUES(20, \"NEW POST EH\")")
	db.MustExec("INSERT INTO posts (owner, text) VALUES(21, \"another one EH\")")
	p := new(Post)
	p.OwnerID = 10
	p.Text = "New method eh"
	p.CreatePost(db)

	p.OwnerID = 10
	p, _ = GetPost(db, p)
	// fmt.Printf("%p\n", p)
	fmt.Println(p)
}
