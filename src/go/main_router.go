package main

import (
	"context"
	"github.com/nvellon/hal"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// PORT is used to configure server port
	PORT = 8080
)

func setupRouter() *gin.Engine {
	db := Connect()
	go DummyUsers(db)
	r := gin.Default()

	// Ping route1¡
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	userDetails := func(c *gin.Context) {
		var u = UserByUsername(c.Params.ByName("username"), db)
		// fmt.Println(user)
		ur := hal.NewResource(u, c.Request.URL.String())
		c.JSON(http.StatusOK, ur)
	}

	authenticated := r.Group("/")
	authenticated.Use(JWTAuth)
	{
		authenticated.GET("/u/:username", userDetails)
	}

	r.POST("/u/:username/login", func(c *gin.Context) {
		u := UserByUsername(c.Params.ByName("username"), db)
		length, err := strconv.Atoi(c.Request.Header.Get("Content-Length"))
		if err != nil {
			log.Fatal(err)
		}
		body := make([]byte, length)
		length, _ = c.Request.Body.Read(body)
		p := ToPassword(body)
		if CheckPass(p, u.PasswordHash) {
			fmt.Println("Same guy")
			jchan := make(chan string)
			go u.GetJWT(&jchan)
			c.JSON(http.StatusOK, Auth{JWT: <-jchan, Valid: true})
		} else {
			fmt.Println("Wrong pass")
			jwt := "Invalid Password"
			c.JSON(http.StatusForbidden, Auth{JWT: jwt, Valid: false})
		}
	})

	return r
}

func runServer(engine *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(PORT),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func main() {
	r := setupRouter()
	runServer(r)
}
