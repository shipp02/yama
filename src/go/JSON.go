package main

// This file contains structs for JSON interactions.

import (
	"github.com/nvellon/hal"
)

// Password is used for authentication
type Password struct {
	Password string
}

type Auth struct {
	JWT   string
	Valid bool
}

func (auth Auth) GetMap() hal.Entry {
	return hal.Entry{
		"jwt":   auth.JWT,
		"valid": auth.Valid,
	}
}

type Response struct {
	Length  int
	Content string
}

func (resp Response) GetMap() hal.Entry {
	return hal.Entry{
		"Length":  resp.Length,
		"Content": resp.Content,
	}
}
