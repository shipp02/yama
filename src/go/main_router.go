package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const(
	// PORT is used to configure server port
	PORT=8080
)



func setupRouter() *gin.Engine {
	db := Connect()
	DummyUsers(db)
	r:= gin.Default()

	// Ping route
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/users/:name",
	func(c *gin.Context){
		user := c.Params.ByName("name")
		fmt.Println(user)
		var u = new(User)
		u.name = user
		nu, _ := GetUser(db, u)
		c.JSON(http.StatusOK, gin.H{"id": nu.id, "username": nu.username})
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("/admin", func(c *gin.Context) {

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	r.POST("/users/:name/login", func(c *gin.Context) {
		name := c.Params.ByName("name")
		var u = new(User)
		u.name = name
		u, _ = GetUser(db, u)
		i:= new(interface{})
		c.ShouldBindBodyWith(*i, binding.JSON)
		fmt.Println("/users/:name/login",*i)
	})

	return r
}

func runServer(engine *gin.Engine){
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(PORT),
		Handler: engine,
	}

	go func(){
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
