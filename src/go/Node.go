package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func newMNode(ID int64) *mNode {
	return &mNode{ID: ID}
}

func (node mNode) GetMap() hal.Entry {
	return hal.Entry{
		"id":           node.ID,
		"name":         node.Name,
		"has_children": node.Children,
	}
}

type Nodes []mNode

func NodeByParent(parentId int64, name string, db *sqlx.DB) (node *mNode) {
	stmt, err := db.Preparex("CALL FindChildDetails(?, ?)")
	if err != nil {
		return
	}
	err = stmt.Get(node, parentId, name)
	if err != nil {
		return
	}
	return
}

func (nodes Nodes) GetMap() hal.Entry {
	entries := make([]hal.Resource, len(nodes))
	entry := hal.Resource{}
	for i, elem := range nodes {
		entries[i] = hal.Resource{
			Payload: elem,
		}
		entry.Embed("", &entries[i])
	}
	return entry.GetMap()
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

func (node *mNode) CreateChild(name string, db *sqlx.DB) error {
	stmt, err := db.Preparex("CALL CreateChild(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(node.ID, name)
	if err != nil {
		return errors.New("creating child failed")
	}
	return nil
}

func (node *mNode) CheckChild(name string, db *sqlx.DB) (possibleName string) {
	stmt, err := db.Preparex("SELECT NodeName(?, ?) AS name")
	if err != nil {
		return ""
	}
	err = stmt.Select(possibleName, node.ID, name)
	if err != nil {
		return ""
	}
	return
}

func (node *mNode) FindChildren(fullPath string, db *sqlx.DB) *Nodes {
	var path = strings.Split(fullPath, "/")[1:]
	currentNode := mNode{}
	stmt, err := db.Preparex("SELECT FindChild(?,  ?) AS \"node.id\"")
	if err != nil {
		log.Println(err)
		return nil
	}
	prev := node.ID
	for _, elem := range path {
		fmt.Println(elem)
		err = stmt.Get(&currentNode, prev, elem)
		if err != nil {
			log.Println(err)
			return nil
		}
		prev = currentNode.ID
	}
	return (*Nodes)(currentNode.GetChildren(1, db))
}

func NodeByID(id int64, db *sqlx.DB) *mNode {
	stmt, err := db.Preparex("SELECT id AS 'node.id', name AS 'node.name', document_id AS 'node.document_id' FROM node WHERE id = ?")
	if err != nil {
		log.Println(err)
		return nil
	}
	var node mNode
	err = stmt.Get(&node, id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &node
}
