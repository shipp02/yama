package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type Group struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type UsersInGroup []mUsers

func GetGroup(id int64, db *sqlx.DB) (*Group, error) {
	stmt, err := db.Preparex("SELECT id AS 'group.id' ,name AS 'group.name' FROM grp WHERE id = ?")
	if err != nil {
		return nil, errors.New("unable to resolve statement")
	}

	var grp Group
	err = stmt.Select(&grp, id)
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
	var users = struct {
		userIds []int64 `db:"user_id"`
	}{}
	err = stmt.Get(&users, grp.ID)
	if err != nil {
		return nil, err
	}
	return users.userIds, nil
}

func (grp *Group) GetUserDetails(db *sqlx.DB) ([]mUsers, error) {
	stmt, err := db.Preparex("SELECT * FROM users_in_groups")
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to resolve statement")
	}
	var ugs []struct {
		userId  int64 `db:"group.user_id"`
		groupId int64 `db:"group.group_id"`
		user    User
	}

	err = stmt.Get(&ugs)
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to resolve statement")
	}
	users := make([]mUsers, len(ugs))
	for i, elem := range ugs {
		users[i] = mUsers{
			ID:       int32(elem.user.Id),
			Name:     elem.user.Name,
			Username: elem.user.Username,
		}
	}
	return users, nil
}
