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
	binHash := sha256.Sum256(pass)
	dst := make([]byte, hex.EncodedLen(len(binHash)))
	hex.Encode(dst, binHash[:])
	return string(dst)
}

func Pbkdf2(s string) (hash string) {
	pass := []byte(s)
	salt := make([]byte, 20)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal("func Pbkdf2 Password failed", err)
	}
	binHash := pbkdf2.Key(pass, salt, 4096, 256, sha256.New)
	strSalt := base64.StdEncoding.EncodeToString(salt)
	strHash := base64.StdEncoding.EncodeToString(binHash)
	return strSalt + ":" + strHash
}

// CheckPass checks the password return true if correct
func CheckPass(pass *Password, hash string) (t bool) {
	// passHash := Sha256(pass.Password)
	// fmt.Println("func CheckPass:hash:", hash)
	passParts := strings.Split(hash, ":")
	passHash, _ := base64.StdEncoding.DecodeString(passParts[1])
	salt, _ := base64.StdEncoding.DecodeString(passParts[0])
	log.Printf(passParts[0])
	binHash := pbkdf2.Key([]byte(pass.Password), salt, 4096, 256, sha256.New)
	if bytes.Equal(binHash, passHash) {
		return true
	}
	return false
}

func (m *mUsers) GetJWT(jwtChan *chan string) {
	conf := DefaultConfig()
	var jc jwt.Claims
	jc.Issuer = *conf.Issuer
	jc.Subject = m.Username
	jc.KeyID = m.PasswordHash
	jwtToken, err := jc.HMACSign(jwt.HS512, *conf.Secret)
	if err != nil {
		fmt.Println("GetUser jwt", err)
	}
	jChan := *jwtChan
	jChan <- base64.StdEncoding.EncodeToString(jwtToken)
	close(jChan)
}

func CheckJWT(strJwt string) bool {
	binJwt, err := base64.StdEncoding.DecodeString(strJwt)
	if err != nil {
		log.Println(err, "func CheckJWT base64 decode failed")
	}
	conf := DefaultConfig()
	jc, err := jwt.HMACCheck(binJwt, *conf.Secret)
	if err != nil || jc == nil {
		log.Println(err)
		return false
	}
	if jc.Issuer == *conf.Issuer {
		return true
	}
	return false
}

func mainC() {
	fmt.Println(Sha256("Aashay"))
}

var JWTAuth = func(c *gin.Context) {
	strJwt := c.Request.Header.Get("X-Auth-Key")
	if CheckJWT(strJwt) {
	} else {
		c.JSON(http.StatusForbidden, gin.H{"Cease": "Desist"})
		c.Abort()
	}
}
