package main

import (
	"context"
	"github.com/nvellon/hal"
	"strings"

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

	authenticated := r.Group("/")
	authenticated.Use(JWTAuth)
	{
		var (
			userDetails gin.HandlerFunc // Gives details of user
			createPost  gin.HandlerFunc // Creates a post
			viewPosts   gin.HandlerFunc // Shows all posts
		)
		authenticated.GET("/u/:username", userDetails)
		authenticated.POST("/edit/u/:username/posts", createPost)
		authenticated.GET("/u/:username/posts", viewPosts)

		userDetails = func(c *gin.Context) {
			var u = UserByUsername(c.Params.ByName("username"), db)
			// fmt.Println(user)
			if u.ID != 0 {
				ur := hal.NewResource(u, c.Request.URL.String())
				c.JSON(http.StatusOK, ur)
			} else {
				s := "user with username %s does not exist"
				c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf(s, c.Params.ByName("username"))})
			}
		}
		createPost = func(c *gin.Context) {
			user := UserByUsername(c.Params.ByName("username"), db)
			var s = struct {
				Text string
			}{}
			err := c.BindJSON(&s)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			post := mPost{
				OwnerID: user.ID,
				Text:    &s.Text,
			}
			err = post.CreatePost(db)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.AbortWithStatus(http.StatusOK)

		}
		viewPosts = func(c *gin.Context) {
			user := UserByUsername(c.Params.ByName("username"), db)
			if user == nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			posts, err := user.GetPosts(db)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			resp := hal.NewResource(Response{Length: len(*posts), Content: "posts"}, c.Request.URL.String())
			p := *posts
			for _, rrp := range p {
				resp.Embedded.Add("posts", hal.NewResource(rrp, ""))
			}
			c.JSON(http.StatusOK, resp)
		}
	}

	{
		var (
		//getChildren gin.HandlerFunc
		//getParents gin.HandlerFunc
		)
		//authenticated.GET("/node/down/:name", getChildren)
		//getChildren = func(c *gin.Context) {
		//	fmt.Println("getChildren", c.Params.ByName("name"))
		//}
	}

	r.POST("/u/:username/login", func(c *gin.Context) {
		u := UserByUsername(c.Params.ByName("username"), db)
		//Length, err := strconv.Atoi(c.Request.Header.Get("Content-Length"))
		var p = new(Password)
		err := c.BindJSON(p)
		if err != nil {
			log.Fatal(err)
		}
		if CheckPass(p, u.PasswordHash) {
			fmt.Println("Same guy")
			jChan := make(chan string)
			go u.GetJWT(&jChan)
			authRes := hal.NewResource(Auth{JWT: <-jChan, Valid: true}, c.Request.RequestURI)
			c.JSON(http.StatusOK, authRes)
		} else {
			fmt.Println("Wrong pass")
			jwt := "Invalid Password"
			authRes := hal.NewResource(Auth{JWT: jwt, Valid: false}, c.Request.RequestURI)
			c.JSON(http.StatusForbidden, authRes)
		}
	})

	createUser := func(c *gin.Context) {
		log.Println(":func createUser ran")
		val := new(mUsers)
		err := c.BindJSON(val)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			log.Fatal(err)
			return
		}
		log.Println("func createUser", val)
		err = val.CreateUser(db)
		if err != nil {
			if strings.Contains(err.Error(), "user exists") {
				c.JSON(http.StatusOK, gin.H{"Error": "User exists"})
				c.Abort()
			} else {
				c.JSON(http.StatusInternalServerError, err)
			}
			return
		}
		jChan := make(chan string)
		go val.GetJWT(&jChan)
		c.JSON(http.StatusOK, Auth{JWT: <-jChan, Valid: true})
	}
	r.POST("/edit/create/u", createUser)

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
