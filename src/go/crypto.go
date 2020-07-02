package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"fmt"

	"github.com/pascaldekloe/jwt"
)

// Sha256 returns hex hash string of a string
func Sha256(s string)(hash string){
	pass := []byte(s)
	binhash := sha256.Sum256(pass)
	dst := make([]byte, hex.EncodedLen(len(binhash)))
	hex.Encode(dst, binhash[:])
	return string(dst)
}

// CheckPass checks the password return true if correct
func CheckPass(pass *Password, hash string)(t bool){
	passHash := Sha256(pass.Password)
	if passHash == hash{
		return true
	} 
	return false
}

func (u *User) GetJWT(jwtchan *chan []byte){
	conf := DefaultConfig()
	var jc jwt.Claims
	jc.Issuer = *conf.Issuer
	jc.Subject = u.Username
	jc.KeyID = u.PasswordHash
	token,err := jc.RSASign(jwt.RS512, conf.PrivateKey)
	btoken := make([]byte, base64.StdEncoding.EncodedLen(len(token)))
	base64.StdEncoding.Encode(btoken, token)
	if err != nil {
		fmt.Println("GetUser jwt",err)
	}
	fmt.Println(string(btoken))
	jchan := *jwtchan
	jchan <-  btoken
}

func mainC(){
	fmt.Println(Sha256("Aashay"))
}