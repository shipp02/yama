package main

// This file contains structs for JSON interactions.

import (
	"github.com/RichardKnop/jsonhal"
	"github.com/nvellon/hal"
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
