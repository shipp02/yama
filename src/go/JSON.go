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
	element hal.Mapper
	Length  int
	Content string
}

func (resp Response) GetMap() hal.Entry {
	if resp.element != nil {
		return hal.Entry{
			"Length":  resp.Length,
			"Content": resp.Content,
			"element": resp.element,
		}
	} else {
		return hal.Entry{
			"Length":  resp.Length,
			"Content": resp.Content,
		}
	}
}

func NodeToMap(nodes *[]mNode) *[]hal.Mapper {
	mapper := make([]hal.Mapper, len(*nodes))
	for i, elem := range *nodes {
		mapper[i] = elem
	}
	return &mapper
}
func UserToMap(users *[]mUsers) *[]hal.Mapper {
	mapper := make([]hal.Mapper, len(*users))
	for i, elem := range *users {
		mapper[i] = elem
	}
	return &mapper
}

func MapArray(objects []hal.Mapper, selfUri string, content string) *hal.Resource {
	resp := Response{
		Length:  len(objects),
		Content: content,
	}
	respResource := hal.NewResource(resp, selfUri)
	for _, object := range objects {
		respResource.Embed("", hal.NewResource(object, ""))
	}
	return respResource
}
