package main

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type mDocument struct {
	Data     sql.RawBytes `db:"document.content"`
	DataType string       `db:"document.type"`
	Name     string
	ID       int64
}

// TODO:Prevent overwriting existing document_id
func (node *mNode) AddDocument(data []byte, docType string, db *sqlx.DB) (*mDocument, error) {
	stmt, err := db.Preparex("INSERT INTO document (name, content, type) VALUES (?, ?, ?)")
	if err != nil {
		return nil, errors.New("unable to resolve stmt")
	}
	exec, err := stmt.Exec(node.Name, data, docType)
	if err != nil {
		return nil, err
	}
	id, _ := exec.LastInsertId()
	doc := mDocument{
		ID:       id,
		Name:     node.Name,
		Data:     data,
		DataType: docType,
	}
	return &doc, nil
}
