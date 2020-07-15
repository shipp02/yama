package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
)

type mNode struct {
	ID         int64         `db:"node.id"`
	Name       string        `db:"node.name"`
	Children   bool          `db:"node.children"`
	ParentID   int64         `db:"node.parent_id"`
	DocumentID sql.NullInt64 `db:"node.document_id"`
	Depth      sql.NullInt64 `db:"node.depth"`
}

func (p *mNode) GetParents(depth int, db *sqlx.DB) *[]mNode {
	stmt, err := db.Preparex("CALL GetParents(?, ?)")
	if err != nil {
		return nil
	}
	var parents []mNode
	err = stmt.Select(&parents, depth, p.ID)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &parents
}

func (n *mNode) GetChildren(depth int, db *sqlx.DB) *[]mNode {
	stmt, err := db.Preparex("CALL GetChildren(?, ?)")
	if err != nil {
		return nil
	}
	var children []mNode
	err = stmt.Select(&children, depth, n.ID)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &children
}
