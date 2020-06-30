package main
// This file contains structs for JSON interactions. 

import (
	"encoding/json"
	"fmt"
)

// Password is used for authentication
type Password struct{
	Password string	
}

func ToPassword(data []byte) (*User){

	user := &User{}
	json.Unmarshal(data, user)

	fmt.Println("JSON:", user)
	return user
}