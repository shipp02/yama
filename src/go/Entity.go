package main;

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	// pass "./password"
)

const (  // iota is reset to 0
	STUDENT = iota  // c0 == 0
	OWNER = iota  // c1 == 1
	ADMIN = iota  // c2 == 2
)


// User Represents user of application
type User struct {
	id int64
	name string
	username string
	passwordHash string
	groups []Group
	permissions []int
}

const userSchema = `
CREATE TABLE users (
	id int NOT NULL AUTO_INCREMENT,
	name VARCHAR(50) NOT NULL ,
	username VARCHAR(100) NOT NULL ,
	password_hash VARCHAR(64) NOT NULL,
	PRIMARY KEY (id)
)`

type queryUser struct {
	id sql.NullInt64
	name sql.NullString
	username sql.NullString
	passwordHash sql.NullString
}

func GetInterface (l int, qu *queryUser) (iface []interface{}){
	iface = make([]interface{}, l)
	iface[0] = &qu.id
	iface[1] = &qu.name
	iface[2] = &qu.username
	iface[3] = &qu.passwordHash
	return
}

func (qu queryUser) ToUser () (u *User){
	fmt.Println("ToUser",qu)
	u = new(User)
	u.id = qu.id.Int64
	u.name = qu.name.String
	u.username = qu.username.String
	u.passwordHash = qu.passwordHash.String
	return
}


type Group struct {
	members []User
	owner User
}

type Folder struct {
	view Group
}

func CleanUp (db *sqlx.DB){
	db.MustExec("USE mysql")
	db.MustExec("DROP TABLE users")
	db.Close()
}

func Connect() (db *sqlx.DB){
	db,err:= sqlx.Connect("mysql", "root:yoursql@tcp(localhost:3306)/mysql")
	if err != nil {
        log.Fatalln(err)
    }
	return 
}

// GetUser filled in from database
func GetUser(db *sqlx.DB, pu *User) (*User, error) {
	var error error
	if pu.id ==0 && pu.name == "" && pu.username == ""{
		error =errors.New("Insufficient data")
	}
	var query = `
	SELECT * FROM users 
	WHERE
	`
	const idQ = "id=$(ID)\n"
	const nameQ = "name=\"$(NAME)\"\n"
	const usernameQ = "username=\"$(UNAME)\"\n"
	if pu.id != 0 {
		IDQ := strings.Replace(idQ, "$(ID)", strconv.FormatInt(pu.id, 10), 1)
		query = query + IDQ
	}
	if pu.name != "" {
		NAMEQ := strings.Replace(nameQ, "$(NAME)", pu.name, 1)
		query = query + NAMEQ
	}
	if pu.username != "" {
		UNAMEQ := strings.Replace(usernameQ, "$(UNAME)", pu.username, 1)
		query = query + UNAMEQ
	}
	fmt.Println(query)
	resp,err := db.Query(query)
	if err != nil{
		log.Fatal("Query Unsatisfied" + query +  "\n" + err.Error())
		error = errors.New("Query Unsatisfied" + query +  "\n" + err.Error())
	}
	l,err := resp.Columns()
	if(err != nil){
		fmt.Println(err)
	}
	var qu queryUser
	var s = GetInterface(len(l), &qu)

	for resp.Next() {
		if err:=resp.Scan(s...); err !=nil {
			log.Fatal(err)
			error = errors.New(err.Error())
		}
		fmt.Println()
		fmt.Println(qu)
		fmt.Println(s...)
		fmt.Println(s)
		fmt.Printf("%p", &qu.id)
		fmt.Println()
		// s... is filled here
	}
	
	if qu.passwordHash.String == ""{
		error = errors.New("Could not find user")
	}
	fmt.Println("query User: ",qu)
	return qu.ToUser(), error
}

// CreateUser creates an entry for User in database
func (u User) CreateUser (db *sqlx.DB) (error){
	var pu  = new(User)
	var error error
	pu.username= u.username
	pu, err := GetUser(db, pu)
	if err == nil{
		error = errors.New("User already exists")
	}
	if u.name == "" || u.username=="" || u.passwordHash==""{
		error = errors.New("User incomplete")
	}
	if error == nil {
		var execu = "INSERT INTO users (username, name, password_hash)VALUES(\"$(UNAME)\", \"$(NAME)\", SHA2(\"$(PASS)\",256))"
		execu =  strings.Replace(execu, "$(UNAME)", u.username, 1)
		execu = strings.Replace(execu, "$(PASS)", u.passwordHash, 1)
		execu = strings.Replace(execu, "$(NAME)", u.name, 1)
		fmt.Println("Create String: "+execu)
		db.MustExec(execu)
	}
	return error
}

// Authenticate checks password against database
func (u User) Authenticate (db *sqlx.DB) (error error) {
	// hash := pass.Sha256(u.passwordHash)
	// pu, err := GetUser(db, &u)
	// if err != nil {
		// log.Println(err.Error())
	// }
	// if pu.passwordHash == hash {
	// 	fmt.Println("Same guy")
	// }
	return error
}

func run() {
	db:= Connect() 
	db.MustExec("USE  mysql")
	db.MustExec(userSchema)

	var pu *User = new(User)
	pu.id = 100
	pu, error := GetUser(db, pu)
	if error != nil{
		fmt.Println(error.Error())
	}
	fmt.Println(pu.username, pu.passwordHash, "User")
	db.Close()
}

func run2() {
	db := Connect() 
	var u User
	u.name= "Chasma"
	u.username = "Devi"
	u.passwordHash = "KALI MA"
	err := u.CreateUser(db)
	fmt.Println(err)
	db.Close()
}

// DummyUsers creates dummy users for use in testing
func DummyUsers(db *sqlx.DB){
	u1:=new(User)
	u1.name = "George"
	u1.username="210978"
	u1.passwordHash="Hkis210978"
	u1.CreateUser(db)
	u2:=new(User)
	u2.name="John"
	u2.username="teacher"
	u2.passwordHash="Yes,papa!"
	u2.CreateUser(db)
}

func main(){
	db:= Connect() 
	run()
	DummyUsers(db)
	run2()
	pu := new(User)
	pu.name = "George"
	user, _:= GetUser(db, pu)
	fmt.Println(user)
	// CleanUp(db)
	db.Close()
}