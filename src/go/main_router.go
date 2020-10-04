package main

import (
	"context"
	"github.com/jmoiron/sqlx"
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

var (
	db     *sqlx.DB
	engine *gin.Engine
)

func init() {
	db = Connect()
	initDocStmt(db)
	engine = setupRouter()
}

func setupRouter() *gin.Engine {
	go DummyUsers(db)
	r := gin.Default()

	// Ping route1¡
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.MaxMultipartMemory = 8 << 20

	authenticated := r.Group("/")
	authenticated.Use(JWTAuth)
	addUserHandling(authenticated, db)
	addPostHandling(authenticated, db)
	addNodeHandling(authenticated, db)
	addGroupHandling(authenticated, db)

	r.POST("/u/:username/login", func(c *gin.Context) {
		u := UserByUsername(c.Params.ByName("username"), db)

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
	r.POST("/edit/u/create", createUser)

	addDocumentHandling(authenticated, db)
	return r
}

func addPostHandling(r *gin.RouterGroup, db *sqlx.DB) {

	createPost := func(c *gin.Context) {
		var s = struct {
			Text string `json:"text"`
		}{}
		err := c.BindJSON(&s)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		post := mPost{
			OwnerID: c.Keys["user"].(mUsers).ID,
			Text:    &s.Text,
		}
		err = post.CreatePost(db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	}

	viewPosts := func(c *gin.Context) {
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
	view1Post := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Coming": "SOON"})
	}
	var edit1Post gin.HandlerFunc
	edit1Post = func(c *gin.Context) {
		//user := c.Keys["user"].(mUsers)
	}
	r.POST("/edit/p/:username", createPost)
	r.PUT("/edit/p/:username/:id", edit1Post)
	r.GET("/view/p/:username", viewPosts)
	r.GET("/view/p/:username/:id", view1Post)
}

func addUserHandling(r *gin.RouterGroup, db *sqlx.DB) {
	userDetails := func(c *gin.Context) {
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
	r.GET("/u/:username", userDetails)

}

func addDocumentHandling(r *gin.RouterGroup, db *sqlx.DB) {
	addDocument := func(c *gin.Context) {
		nodeId, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ":id must be int"})
			return
		}
		node := NodeByID(int64(nodeId), db)
		formFile, err := c.FormFile("file")
		if formFile == nil || err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		file, err := formFile.Open()
		if file == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		dataSpace := make([]byte, 8<<20)
		read, err := file.Read(dataSpace)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		data := dataSpace[:read]
		document, err := node.AddDocument(data, "pdf", db)
		c.JSON(http.StatusOK, document)
	}
	r.POST("/edit/d/:id", addDocument)
	viewDocument := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ":id must be integer"})
			return
		}
		document, err := DocumentByID(int64(id), db)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp := hal.NewResource(document, c.Request.RequestURI)
		c.JSON(http.StatusOK, resp)
	}
	r.GET("/view/d/:id", viewDocument)

}

func addNodeHandling(r *gin.RouterGroup, db *sqlx.DB) {
	createChild := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "id must be a integer"})
		}
		parentNode := newMNode(int64(id))
		var name = struct {
			Name string `json:"name"`
		}{}
		err = c.BindJSON(&name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "name not provided in json"})
			log.Println(err.Error())
			return
		}
		err = parentNode.CreateChild(name.Name, db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		node := NodeByParent(int64(id), name.Name, db)
		//resp := hal.NewResource(*node, c.Request.RequestURI)
		c.JSON(http.StatusOK, *node)
	}
	r.POST("/edit/tree/:id", createChild)
	getChildren := func(c *gin.Context) {
		fullPath := c.Params.ByName("name")
		node := mNode{
			ID: 1,
		}
		//res := hal.NewResource(node.FindChildren(fullPath, db), c.Request.RequestURI)
		nodes := node.FindChildren(fullPath, db)
		//resp := MapArray(*NodeToMap(nodes), c.Request.RequestURI, "children")
		//resp := hal.NewResource(nodes, c.Request.RequestURI)
		//resp := NodeToMap((*[]mNode)(nodes))
		c.JSON(http.StatusOK, nodes)
	}
	r.GET("/view/tree/down/~/*name", getChildren)
	getChildrenContext := func(c *gin.Context) {
		contextPath := c.Request.Header.Get("Context")
		fullPath := contextPath + c.Params.ByName("path")
		node := mNode{
			ID: 1,
		}
		nodes := *node.FindChildren(fullPath, db)
		c.JSON(http.StatusOK, nodes)
	}
	r.GET("/view/tree/down/_/*path", getChildrenContext)
	checkChild := func(c *gin.Context) {}
	r.POST("/edit/tree/:id/check", checkChild)

}

func addGroupHandling(r *gin.RouterGroup, db *sqlx.DB) {

	viewGroups := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": ":id must be a valid integer"})
			return
		}
		grp, err := GetGroup(id, db)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "requested group does not exist"})
			return
		}
		members, err := grp.GetUserDetails(db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var resp = hal.NewResource(Response{Content: "members", Length: len(members), element: grp}, c.Request.RequestURI)
		for _, member := range members {
			resp.Embedded.Add(hal.Relation(strconv.Itoa(int(member.ID))), hal.NewResource(member, "/u/"+member.Username))
		}
		c.JSON(http.StatusOK, resp)
	}
	createGroup := func(c *gin.Context) {
		name := c.Params.ByName("name")
		if name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "name field cannot be empty"})
		}
		grp := Group{Name: name}
		err := grp.CreateGroup(db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		resp := hal.NewResource(grp, "/view/g/"+strconv.Itoa(int(grp.ID)))
		c.JSON(http.StatusOK, resp)
	}
	addUserToGrp := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ":id  must be an integer"})
			return
		}
		var u mUsers
		err = c.BindJSON(&u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body must contain valid json"})
			return
		}
		err = u.AddToGroup(int64(id), db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user or group does not exist"})
			return
		}
		c.Status(http.StatusOK)
	}
	r.GET("/view/g/:id", viewGroups)
	r.GET("/edit/g/:name/create", createGroup)
	r.PUT("/edit/g/:id/add", addUserToGrp)
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
	runServer(engine)
}
