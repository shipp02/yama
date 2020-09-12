package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
)

var c = DefaultConfig()

type Config struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Issuer     string
	Secret     []byte
}

func DefaultConfig() *Config {
	fmt.Println("Create config")
	c := new(Config)
	var err error
	c.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	c.PublicKey = &c.PrivateKey.PublicKey
	c.Issuer = "yama"
	c.Secret = []byte("Sec key")
	return c
}
