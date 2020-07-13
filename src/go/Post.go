package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"../resources/test/model"
	. "../resources/test/table"
	. "github.com/go-jet/jet/v2/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// PostSchema Schema for posts table
const PostSchema = `
CREATE TABLE posts(
	id int NOT NULL AUTO_INCREMENT,
	owner_id int NOT NULL,
	text TEXT,
	PRIMARY KEY(id),
	FOREIGN KEY (owner_id) REFERENCES users (id) 
);
`

// Post represents a post on a message board
type Post struct {
	ID      int64  `db:"id"`
	OwnerID int64  `db:"owner_id"`
	Text    string `db:"text"`
}

type mPost model.Posts

// GetPost will fetch particular post from db
func GetPost(db *sqlx.DB, p *mPost) (*mPost, error) {
	var err error
	if p.ID == 0 && p.OwnerID == 0 {
		err = errors.New("insufficient data")
	}
	stmt := SELECT(Posts.ID.AS("mPosts.id"),
		Posts.OwnerID.AS("mPosts.owner_id"),
		Posts.Text.AS("mPosts.Text")).FROM(Posts).
		WHERE(Posts.ID.EQ(Int(int64(p.ID)))).
		LIMIT(1)
	err = stmt.Query(db, p)
	if err != nil {
		if strings.Contains(err.Error(), "qrm: no rows in result set") {
			return p, nil
		} else {
			return nil, err
		}
	}
	// fmt.Printf("%p\n", p)
	return p, err
}

// GetPosts gets all posts of a user
func (u *User) GetPosts(db *sqlx.DB) ([]Post, error) {
	var err error
	posts := []Post{}
	query := `
		SELECT * FROM posts
		WHERE owner_id=$(OID)
	`
	query = strings.Replace(query, "$(OID)", strconv.FormatInt(u.Id, 10), 1)
	// fmt.Println(query)
	err2 := db.Select(&posts, query)
	if err2 != nil {
		fmt.Println(err2)
	}
	// fmt.Println(posts)
	return posts, err
}

// CreatePost Stores the post in database
func (p *mPost) CreatePost(db *sqlx.DB) error {
	qp, err := GetPost(db, p)
	if qp.ID != 0 && qp != nil {
		err = errors.New("post exists")
	}

	jetFlag := true
	if jetFlag {
		stmt := Posts.INSERT(Posts.OwnerID, Posts.Text).VALUES(p.OwnerID, p.Text)
		_, err := stmt.Exec(db)
		if err != nil {
			return err
		}

	} else {
		var exec = "INSERT INTO posts (owner_id, text) VALUES($(OID), \"$(TEXT)\")"

		exec = strings.Replace(exec, "$(OID)", strconv.FormatInt(int64(p.OwnerID), 10), 1)
		exec = strings.Replace(exec, "$(TEXT)", *p.Text, 1)
		// fmt.Println(exec)

		db.MustExec(exec)
	}
	return err
}

func mainP() {
	db := Connect()
	db.MustExec(PostSchema)
	db.MustExec("INSERT INTO posts (owner_id, text) VALUES(20, \"NEW POST EH\")")
	db.MustExec("INSERT INTO posts (owner_id, text) VALUES(21, \"another one EH\")")
	p := new(mPost)
	p.OwnerID = 10
	s := "New method eh"
	p.Text = &s
	_ = p.CreatePost(db)

	p.OwnerID = 10
	p, _ = GetPost(db, p)
	// fmt.Printf("%p\n", p)
	fmt.Println(p)
}
