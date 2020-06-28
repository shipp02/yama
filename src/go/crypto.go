package main;

import (
	"fmt"
	"encoding/hex"
	"crypto/sha256"
)

// Sha256 returns hex hash string of a string
func Sha256(s string)(hash string){
	pass := []byte(s)
	binhash := sha256.Sum256(pass)
	dst := make([]byte, hex.EncodedLen(len(binhash)))
	hex.Encode(dst, binhash[:])
	return string(dst)
}

func CheckPass(pass string, hash string)(t bool){
	passHash := Sha256(pass)
	if passHash == hash{
		return true
	} 
	return false
}

func mainC(){
	fmt.Println(Sha256("Aashay"))
}