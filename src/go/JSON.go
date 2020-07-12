package main

// This file contains structs for JSON interactions.

import (
	"encoding/json"

	"github.com/RichardKnop/jsonhal"
)

// Password is used for authentication
type Password struct {
	Password string
}

type Auth struct {
	jsonhal.Hal
	JWT   string
	Valid bool
}

// func Auth(JWT *string, Valid bool)(a Auth){
// 	a = Auth{*JWT, Valid}
// 	return
// }

func ToPassword(data []byte) (pass *Password) {
	pass = new(Password)
	err := json.Unmarshal(data, pass)
	if err != nil {
		return nil
	}
	return
}

type UDetails struct {
	jsonhal.Hal
	ID       int64
	Name     string
	Username string
}

func (u *User) ToUDetails() *UDetails {
	ud := UDetails{ID: u.Id, Name: u.Name, Username: u.Username}
	return &ud
}
