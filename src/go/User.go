package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nvellon/hal"
	"log"
	"strconv"
	"strings"

	"../resources/test/model"
	. "../resources/test/table"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	// pass "./password"
	. "github.com/go-jet/jet/v2/mysql"
)

var jetFlag = false

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

const userSchema = `
CREATE TABLE users (
	id int NOT NULL AUTO_INCREMENT,
	name VARCHAR(50) NOT NULL ,
	username VARCHAR(100) NOT NULL ,
	password_hash VARCHAR(373) NOT NULL,
	PRIMARY KEY (id)
)`

type queryUser struct {
	id           sql.NullInt64
	name         sql.NullString
	username     sql.NullString
	passwordHash sql.NullString
}

func (qu *queryUser) GetInterface(l int) (iface []interface{}) {
	iface = make([]interface{}, l)
	iface[0] = &qu.id
	iface[1] = &qu.name
	iface[2] = &qu.username
	iface[3] = &qu.passwordHash
	return
}

func (qu *queryUser) ToUser() (u *User) {
	u = new(User)
	u.Id = qu.id.Int64
	u.Name = qu.name.String
	u.Username = qu.username.String
	u.PasswordHash = qu.passwordHash.String
	return
}

// CleanUp removes tables after tests
func CleanUp(db *sqlx.DB) {
	db.MustExec("DROP TABLE users")
	_ = db.Close()
}

// Connect connects to my database
func Connect() (db *sqlx.DB) {
	db, err := sqlx.Connect("mysql", "root:yoursql@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatalln(err)
	}
	return
}

// GetUser filled in from database
func GetUser(db *sqlx.DB, pu *User) (*User, error) {
	var err2 error
	jetFlag := true
	if pu.Id == 0 && pu.Username == "" {
		err2 = errors.New("Insufficient data")
		return pu, err2
	}
	if jetFlag {
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

	} else {
		var query = `
	SELECT id, name, username, password_hash
	FROM users 
	WHERE
	`
		const idQ = "id=$(ID)\n"
		const nameQ = "name=\"$(NAME)\"\n"
		const usernameQ = "username=\"$(UNAME)\"\n"
		var where string
		if pu.Id != 0 {
			where = strings.Replace(idQ, "$(ID)", strconv.FormatInt(pu.Id, 10), 1)
		}
		if pu.Name != "" && where == "" {
			where = strings.Replace(nameQ, "$(NAME)", pu.Name, 1)
		}
		if pu.Username != "" && where == "" {
			where = strings.Replace(usernameQ, "$(UNAME)", pu.Username, 1)
		}
		resp, err := db.Query(query + where)
		if err != nil {
			log.Println(err)
		}
		l, err := resp.Columns()
		if err != nil {
			fmt.Println(err)
		}
		var qu *queryUser = new(queryUser)
		var s = qu.GetInterface(len(l))

		for resp.Next() {
			if err := resp.Scan(s...); err != nil {
				log.Println(err)
			}
		}

		if qu.passwordHash.String == "" {
			err2 = errors.New("Could not find user")
		}
		log.Println("non jet ran")
		return qu.ToUser(), err2
	}
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
func (u User) CreateUser(db *sqlx.DB) error {
	var err error
	jetFlag := true
	if u.Name == "" || u.Username == "" || u.PasswordHash == "" {
		err = errors.New("User incomplete")
		return err
	}
	if uCheck := UserByUsername(u.Username, db); uCheck.ID != 0 {
		return errors.New("user exists")
	}
	if jetFlag {
		exec := Users.INSERT(
			Users.Name,
			Users.Username,
			Users.PasswordHash,
		).VALUES(
			u.Name,
			u.Username,
			Pbkdf2(u.PasswordHash))
		result, err := exec.Exec(db)
		if err != nil {
			log.Fatal("Create User: ", err)
		}
		log.Println("Create User", result)
	} else {
		var pu = new(User)
		var err2 error
		pu.Username = u.Username
		if err2 == nil {
			//var execu = "INSERT INTO users (username, name, password_hash) VALUES(\"$(UNAME)\", \"$(NAME)\", SHA2(\"$(PASS)\",256))"
			var exec = "INSERT INTO users (username, name, password_hash) VALUES(\"%s\", \"%s\", \"%s\")"
			exec = fmt.Sprintf(exec, u.Username, u.Name, Pbkdf2(u.PasswordHash))
			//execu = strings.Replace(execu, "$(UNAME)", u.Username, 1)
			//execu = strings.Replace(execu, "$(PASS)", u.PasswordHash, 1)
			//execu = strings.Replace(execu, "$(NAME)", u.Name, 1)
			db.MustExec(exec)
		}
		return err2
	}
	return err
}

// DummyUsers creates dummy users for use in testing
func DummyUsers(db *sqlx.DB) {
	//db.MustExec(userSchema)
	u1 := new(User)
	u1.Name = "George"
	u1.Username = "210978"
	u1.PasswordHash = "Hkis210978"
	err := u1.CreateUser(db)
	if err != nil {
		log.Fatal(err)
	}
	u2 := new(User)
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

	u1, err = GetUser(db, u1)
	_, err = u1.GetPosts(db)
	if err != nil {
		fmt.Println(err)
	}
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
