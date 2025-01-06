package main

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/shaj13/libcache"
	// _ "github.com/shaj13/libcache/fifo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

// const creditScoreMin = 500
// const creditScoreMax = 900

// type Task struct {
// 	Id          string `json:"id"`
// 	Task        string `json:"task"`
// 	IsCompleted string `json:"iscompleted"`
// }

// type Result struct {
// 	Value []Task
// 	Err   error
// }

// Global variables

var db *sql.DB
var Rsecret string
var Asecret string

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "postgres"
	dbname   = "Forum"
)

// dynamically bind the data to a struct
func parser(c *gin.Context, data interface{}, message string) {
	// context has already been consumed at this point so we get the raw body
	rawBody, _ := c.Get("rawBody")
	body := rawBody.([]byte)

	// bind the data
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
}

func postHelper(c *gin.Context, message string, sql string, data ...interface{}) {

	// println("CHECKING ID", Creator_Id)
	_, err := db.Exec(sql, data...)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to insert " + message})
		return
	}
	// Respond with success message or status
	c.JSON(200, gin.H{"message": message + " created successfully"})

}

func createSubPost(c *gin.Context) {
	var subPost struct {
		Creator_id  string `json:"creator_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	parser(c, &subPost, "Failed to Bind a subPost")

	Name := subPost.Name
	Description := subPost.Description
	Creator_Id := subPost.Creator_id
	sql := "INSERT INTO subposts (name, description, creator_id) VALUES ($1, $2, $3)"

	postHelper(c, "SubPost", sql, Name, Description, Creator_Id)

}

func createPost(c *gin.Context) {
	var Post struct {
		User_id     string `json:"user_id"`
		SubPostId   string `json:"sub_post_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PostType    string `json:"post_type"`
	}
	// bind the data to the struct
	parser(c, &Post, "Failed to Bind a post")

	User_Id := Post.User_id
	SubPostId := Post.SubPostId
	Title := Post.Title
	Description := Post.Description
	PostType := Post.PostType

	sql := "INSERT INTO posts (user_id, sub_post_id, title,  description, post_type) VALUES ($1, $2, $3, $4, $5)"

	// update the database
	postHelper(c, "Post", sql, User_Id, SubPostId, Description, Title, PostType)

}

func createReply(c *gin.Context) {
	log.Println("SHOULD BE HERE")
	sql := "INSERT INTO replies (user_id, post_id, reply, parent_reply_id) VALUES ($1, $2, $3, $4)"

	var Reply struct {
		User_id       string `json:"user_id"`
		PostId        string `json:"post_id"`
		Replys        string `json:"reply"`
		ParentReplyId string `json:"parent_reply_id"`
	}
	// bind the data to the struct
	parser(c, &Reply, "Failed to Bind a reply")

	User_Id := Reply.User_id
	PostId := Reply.PostId
	Replys := Reply.Replys

	// check if it's an int

	_, err := strconv.Atoi(Reply.ParentReplyId)
	if err == nil {
		// add parent reply id if we have it
		ParentReplyId := Reply.ParentReplyId
		//update the database
		postHelper(c, "Reply", sql, User_Id, PostId, Replys, ParentReplyId)
	} else {
		// if we dont have parent reply id we need to change the sql accordigly
		sql := "INSERT INTO replies (user_id, post_id, reply) VALUES ($1, $2, $3)"
		// update the database
		postHelper(c, "Reply", sql, User_Id, PostId, Replys)
	}
}

func init() {

	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func main() {

	// env := &Env{db: db}

	Myauth.SetupGoGuardian()

	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/v1/subpost/create", Myauth.Middleware(createSubPost))
	router.POST("/v1/post/create", Myauth.Middleware(createPost))
	router.POST("/v1/replies/create", Myauth.Middleware(createReply))
	router.POST("/v1/auth/login", func(c *gin.Context) {

		c.Set("request", c.Request)
		Myauth.LoginHandler(c, db)
	})
	router.POST("/v1/auth/refresh", Myauth.RefreshHandler)
	router.POST("v1/signup", Myauth.Signup)

	// router.GET("/v1/auth/token", middleware(createToken))

	// Start and run the server
	router.Run(":5000")

	fmt.Println("Server is running on http://localhost:5000")
	fmt.Println("Successfully connected!")
}
