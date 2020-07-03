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

type Auth struct {
	JWT string
	Valid bool
}

// func Auth(JWT *string, Valid bool)(a Auth){
// 	a = Auth{*JWT, Valid}
// 	return
// }

func ToPassword(data []byte) (pass *Password){

	pass = new(Password)
	json.Unmarshal(data, pass)

	fmt.Println("JSON:", pass)
	return 
}