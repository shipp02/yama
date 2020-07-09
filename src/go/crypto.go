package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/pbkdf2"
)

// Sha256 returns hex hash string of a string
func Sha256(s string) (hash string) {
	pass := []byte(s)
	binhash := sha256.Sum256(pass)
	dst := make([]byte, hex.EncodedLen(len(binhash)))
	hex.Encode(dst, binhash[:])
	return string(dst)
}

func Pbkdf2(s string) (hash string) {
	pass := []byte(s)
	salt := make([]byte, 20)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal("func Pbkdf2 Password failed", err)
	}
	binhash := pbkdf2.Key(pass, salt, 4096, 256, sha256.New)
	strsalt := base64.StdEncoding.EncodeToString(salt)
	strhash := base64.StdEncoding.EncodeToString(binhash)
	return strsalt + ":" + strhash
}

// CheckPass checks the password return true if correct
func CheckPass(pass *Password, hash string) (t bool) {
	// passHash := Sha256(pass.Password)
	// fmt.Println("func CheckPass:hash:", hash)
	passParts := strings.Split(hash, ":")
	passhash, _ := base64.StdEncoding.DecodeString(passParts[1])
	salt, _ := base64.StdEncoding.DecodeString(passParts[0])
	log.Printf(passParts[0])
	binhash := pbkdf2.Key([]byte(pass.Password), salt, 4096, 256, sha256.New)
	if bytes.Equal(binhash, passhash) {
		return true
	}
	return false
}

func (u *User) GetJWT(jwtchan *chan string) {
	conf := DefaultConfig()
	var jc jwt.Claims
	jc.Issuer = *conf.Issuer
	jc.Subject = u.Username
	jc.KeyID = u.PasswordHash
	jwtToken, err := jc.HMACSign(jwt.HS512, *conf.Secret)
	if err != nil {
		fmt.Println("GetUser jwt", err)
	}
	jchan := *jwtchan
	jchan <- base64.StdEncoding.EncodeToString(jwtToken)
	close(jchan)
}

func CheckJWT(strjwt string) bool {
	binjwt, err := base64.StdEncoding.DecodeString(strjwt)
	if err != nil {
		log.Println(err, "func CheckJWT base64 decode failed")
	}
	conf := DefaultConfig()
	jc, err := jwt.HMACCheck(binjwt, *conf.Secret)
	if err != nil || jc == nil {
		log.Println(err)
		return false
	}
	if jc.Issuer == *conf.Issuer {
		return true
		log.Println("Correct")
	}
	return false
}

func mainC() {
	fmt.Println(Sha256("Aashay"))
}

var JWTAuth = func(c *gin.Context) {
	strjwt := c.Request.Header.Get("X-Auth-Key")
	if CheckJWT(strjwt) {
	} else {
		c.JSON(http.StatusForbidden, gin.H{"Cease": "Desist"})
		c.Abort()
	}
}
