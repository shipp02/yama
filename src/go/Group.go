package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/nvellon/hal"
	"log"
)

var grpStmt grpStmtStruct

type Group struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (grp Group) GetMap() hal.Entry {
	return hal.Entry{
		"id":   grp.ID,
		"name": grp.Name,
	}
}

type grpStmtStruct struct {
	Get        sqlx.Stmt
	Create     sqlx.Stmt
	GetUserIDs sqlx.Stmt
	GetUsers   sqlx.Stmt
}

//func initGrpStmt () grpStmtStruct{
//	grpStmt = grpStmtStruct{
//		Get:        nil,
//		Create:     nil,
//		GetUserIDs: nil,
//		GetUsers:   nil,
//	}
//	return grpStmt
//}

func GetGroup(id int, db *sqlx.DB) (*Group, error) {
	stmt, err := db.Preparex("SELECT id ,name  FROM grp WHERE id = ?")
	if err != nil {
		return nil, errors.New("unable to resolve statement")
	}

	var grp Group
	err = stmt.Get(&grp, id)
	if err != nil {
		return nil, err
	}
	return &grp, nil
}

func (grp *Group) CreateGroup(db *sqlx.DB) error {
	stmt, err := db.Preparex("INSERT INTO grp (name) VALUES ( ?)")
	if err != nil {
		return errors.New("unable to resolve statement")
	}
	result, err := stmt.Exec(grp.Name)
	if err != nil {
		log.Println(err)
		return errors.New("error while inserting")
	}

	grp.ID, _ = result.LastInsertId()
	return nil
}

func (usr *mUsers) AddToGroup(groupId int64, db *sqlx.DB) error {
	stmt, err := db.Preparex("INSERT INTO usergroups (group_id, user_id) VALUES (?,?)")
	if err != nil {
		return errors.New("unable to resolve statement")
	}
	_, err = stmt.Exec(groupId, usr.ID)

	if err != nil {
		log.Println(err)
		return errors.New("error while inserting")
	}
	return nil
}

func (grp *Group) GetUsersID(db *sqlx.DB) ([]int64, error) {
	stmt, err := db.Preparex("SELECT user_id FROM usergroups WHERE group_id = ?")
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to resolve statement")
	}
	var users []int64
	err = stmt.Select(&users, grp.ID)
	if err != nil {
		return nil, err
	}
	//var userIds = make([]int64, len(users))
	return users, nil
}

func (grp *Group) GetUserDetails(db *sqlx.DB) ([]mUsers, error) {
	stmt, err := db.Preparex("SELECT * FROM users_in_groups WHERE `group.group_id` = ?")
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to resolve statement")
	}
	var ugs []struct {
		UserId  int64 `db:"group.user_id"`
		GroupId int64 `db:"group.group_id"`
		User
	}

	err = stmt.Select(&ugs, grp.ID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to resolve statement")
	}
	users := make([]mUsers, len(ugs))
	for i, elem := range ugs {
		users[i] = mUsers{
			ID:       int32(elem.User.Id),
			Name:     elem.User.Name,
			Username: elem.User.Username,
		}
	}
	return users, nil
}
