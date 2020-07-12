package main

import (
	"errors"
	"github.com/nvellon/hal"
	"log"
	"strings"

	"../resources/test/model"
	. "../resources/test/table"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	// pass "./password"
	. "github.com/go-jet/jet/v2/mysql"
)

//var jetFlag = false

// User Represents user of application
type User struct {
	Id           int64
	Name         string
	Username     string
	PasswordHash string
	JWT          string
	// groups []Group
	permissions []int
}

type mUsers model.Users

func (m mUsers) GetMap() hal.Entry {
	return hal.Entry{
		"id":       m.ID,
		"name":     m.Name,
		"username": m.Username,
	}
}

// CleanUp removes tables after tests
func CleanUp(db *sqlx.DB) {
	db.MustExec("DROP TABLE users")
	_ = db.Close()
}

// Connect connects to my database
func Connect() (db *sqlx.DB) {
	db, err := sqlx.Connect("mysql", "root:yousql@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatalln(err)
	}
	return
}

// GetUser filled in from database
func GetUser(db *sqlx.DB, pu *User) (*User, error) {
	var err2 error
	if pu.Id == 0 && pu.Username == "" {
		err2 = errors.New("insufficient data")
		return pu, err2
	}
	var stmt SelectStatement
	if pu.Id != 0 {
		stmt = SELECT(Users.AllColumns).WHERE(
			Users.ID.EQ(Int(pu.Id)),
		).FROM(Users).LIMIT(1)
	} else {
		stmt = SELECT(Users.AllColumns).
			FROM(Users).
			WHERE(Users.Username.EQ(String(pu.Username)))
	}
	dest := new(model.Users)
	err := stmt.Query(db, dest)
	if err != nil {
		log.Println("func GetUser:", err)
		log.Println(stmt.DebugSql())
		return pu, err
	}
	pu.Name = dest.Name
	pu.PasswordHash = dest.PasswordHash
	pu.Username = dest.Username
	log.Println("Jet ran")

	return pu, err2
}

func UserByID(id int64, db *sqlx.DB) *mUsers {
	stmt := SELECT(Users.ID.AS("mUsers.id"),
		Users.Name.AS("mUsers.name"),
		Users.Username.AS("mUsers.username"),
		Users.PasswordHash.AS("mUsers.PasswordHash")).
		FROM(Users).
		WHERE(Users.ID.EQ(Int(id)))

	dest := new(mUsers)
	err := stmt.Query(db, dest)
	if err != nil {
		log.Println("func UserByID:", err)
		log.Println(stmt.DebugSql())
		return dest
	}
	return dest
}

func UserByUsername(username string, db *sqlx.DB) *mUsers {
	stmt := SELECT(Users.ID.AS("mUsers.id"),
		Users.Name.AS("mUsers.name"),
		Users.Username.AS("mUsers.username"),
		Users.PasswordHash.AS("mUsers.PasswordHash")).
		FROM(Users).
		WHERE(Users.Username.EQ(String(username)))

	dest := new(mUsers)
	err := stmt.Query(db, dest)
	if err != nil {
		if strings.Contains(err.Error(), "qrm: no rows in result set") {
			return dest
		} else {
			log.Println("func UserByUsername:", err)
			log.Println(stmt.DebugSql())
			return dest
		}
	}
	log.Println(*dest)
	return dest
}

// CreateUser creates an entry for User in database
func (m *mUsers) CreateUser(db *sqlx.DB) error {
	var err error
	if m.Name == "" || m.Username == "" || m.PasswordHash == "" {
		err = errors.New("user incomplete")
		return err
	}
	if uCheck := UserByUsername(m.Username, db); uCheck.ID != 0 {
		return errors.New("user exists")
	}
	exec := Users.INSERT(
		Users.Name,
		Users.Username,
		Users.PasswordHash,
	).VALUES(
		m.Name,
		m.Username,
		Pbkdf2(m.PasswordHash))
	log.Println(exec.DebugSql())
	result, err := exec.Exec(db)
	if err != nil {
		log.Fatal("Create User: ", err)
	}
	log.Println("Create User", result)
	return err
}

// DummyUsers creates dummy users for use in testing
func DummyUsers(db *sqlx.DB) {
	//db.MustExec(userSchema)
	u1 := new(mUsers)
	u1.Name = "George"
	u1.Username = "210978"
	u1.PasswordHash = "hkis210978"
	err := u1.CreateUser(db)
	if err != nil {
		log.Fatal(err)
	}
	u2 := new(mUsers)
	u2.Name = "John"
	u2.Username = "teacher"
	u2.PasswordHash = "Yes,papa!"
	err = u2.CreateUser(db)
	if err != nil {
		log.Fatal(err)
	}

	//db.MustExec(PostSchema)
	p := new(Post)
	p.OwnerID = 1
	p.Text = "George posts"
	_ = p.CreatePost(db)
	_ = p.CreatePost(db)

	//_, err = u1.GetPosts(db)
	//db.MustExec(NodeSchema)
	//n := new(Node)
	//n.Name.String = "ROOT"
	//n.ParentID.Int64 = -1
	//n.Children.Bool = true
	//err = n.CreateNode(db)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//nq := new(Node)
	//nq.ID.Int64 = 1
	//nq, err = nq.GetNode(db)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// db.MustExec(GroupSchema)
	// db.MustExec(GroupUserSchema)
}

func mainE() {
	db := Connect()
	//run()
	//DummyUsers(db)
	//run2()
	//pu := new(User)
	//pu.Name = "George"
	//user, _ := GetUser(db, pu)
	//fmt.Println(user)
	// CleanUp(db)
	_ = db.Close()
}
