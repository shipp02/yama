package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/nvellon/hal"
	"log"
	"strings"
)

type mNode struct {
	ID         int64         `db:"node.id"`
	Name       string        `db:"node.name"`
	Children   bool          `db:"node.children"`
	ParentID   int64         `db:"node.parent_id"`
	DocumentID sql.NullInt64 `db:"node.document_id"`
	Depth      sql.NullInt64 `db:"node.depth"`
}

func (node mNode) GetMap() hal.Entry {
	return hal.Entry{
		"id":           node.ID,
		"name":         node.Name,
		"has_children": node.Children,
	}
}

func (node *mNode) GetParents(depth int, db *sqlx.DB) *[]mNode {
	stmt, err := db.Preparex("CALL GetParents(?, ?)")
	if err != nil {
		return nil
	}
	var parents []mNode
	err = stmt.Select(&parents, depth, node.ID)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &parents
}

func (node *mNode) GetChildren(depth int, db *sqlx.DB) *[]mNode {
	stmt, err := db.Preparex("CALL GetChildren(?, ?)")
	if err != nil {
		return nil
	}
	var children []mNode
	err = stmt.Select(&children, depth, node.ID)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &children
}

func (node *mNode) CreateChild(name string, db *sqlx.DB) {
	stmt, err := db.Preparex("CALL CreateChild(?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(node.ID, name)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func (node *mNode) FindChildren(fullPath string, db *sqlx.DB) *[]mNode {
	var path = strings.Split(fullPath, "/")[1:]
	currentNode := mNode{}
	stmt, err := db.Preparex("SELECT FindChild(?,  ?) AS \"node.id\"")
	if err != nil {
		log.Println(err)
		return nil
	}
	prev := node.ID
	for _, elem := range path {
		err = stmt.Get(&currentNode, prev, elem)
		if err != nil {
			log.Println(err)
			return nil
		}
		prev = currentNode.ID
	}
	return currentNode.GetChildren(1, db)
}
