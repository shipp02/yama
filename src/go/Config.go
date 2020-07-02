package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"fmt"
)

var c *Config

type Config struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	Issuer *string
}

func DefaultConfig() (*Config) {
	if c == nil {
		c = new(Config)
		var err error
		c.PrivateKey,err = rsa.GenerateKey(rand.Reader, 2048)
		fmt.Println("Private Key:",c.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		c.PublicKey = &c.PrivateKey.PublicKey
		fmt.Println("Public Key:",c.PublicKey)
		iss := "yama"
		c.Issuer = &iss
	}
	return c
}