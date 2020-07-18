package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
)

var c *Config

type Config struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Issuer     *string
	Secret     *[]byte
}

func DefaultConfig() *Config {
	if c == nil {
		fmt.Println("Create config")
		c = new(Config)
		var err error
		c.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		// fmt.Println("Private Key:",c.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		c.PublicKey = &c.PrivateKey.PublicKey
		// fmt.Println("Public Key:",c.PublicKey)
		iss := "yama"
		c.Issuer = &iss
		x := []byte("Sec key")
		c.Secret = &x
	}
	return c
}
