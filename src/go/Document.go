package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/nvellon/hal"
	"log"
)

type MDocument struct {
	Data     []byte `db:"document.content"`
	DataType string `db:"document.type"`
	Name     string `db:"document.name"`
	ID       int64  `db:"document.id"`
}

func (doc MDocument) GetMap() hal.Entry {
	return hal.Entry{
		"data": doc.Data,
		"type": doc.DataType,
		"Name": doc.Name,
		"id":   doc.ID,
	}
}

// TODO:Prevent overwriting existing document_id
func (node *mNode) AddDocument(data []byte, docType string, db *sqlx.DB) (MDocument, error) {
	var doc MDocument
	if !node.DocumentID.Valid {
		doc, _ = node.GetDocument(db)
		return doc, errors.New("document exists")
	}
	stmt, err := db.Preparex("INSERT INTO document (name, content, type) VALUES (?, ?, ?)")
	if err != nil {
		return doc, errors.New("unable to resolve stmt")
	}
	exec, err := stmt.Exec(node.Name, data, docType)
	if err != nil {
		return doc, err
	}
	id, _ := exec.LastInsertId()
	doc = MDocument{
		ID:       id,
		Name:     node.Name,
		Data:     data,
		DataType: docType,
	}
	return doc, nil
}

func (node *mNode) GetDocument(db *sqlx.DB) (MDocument, error) {
	return DocumentByID(node.ID, db)
}

func DocumentByID(id int64, db *sqlx.DB) (MDocument, error) {
	var doc MDocument
	stmt, err := db.Preparex(`SELECT
       id AS 'document.id', 
       name AS 'document.name', 
       content AS 'document.content',
       type AS 'document.type'
	FROM document
	WHERE id = ?`)
	if err != nil {
		log.Println(err)
		return doc, errors.New("unable to create statement")
	}
	err = stmt.Get(&doc, id)
	if err != nil {
		log.Println(err)
		return doc, err
	}
	return doc, nil
}
