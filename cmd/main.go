package main

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/internal/auth"
	"fmt"
	"log"
	"net/http"

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
var jwtSecret []byte     // Secret for access token
var refreshSecret []byte // Secret for refresh token

var db *sql.DB
var Rsecret string
var Asecret string

type Env struct {
	db     *sql.DB
	logger *log.Logger
}

// var tokenStrategy auth.Strategy

// JWT expiration times

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "postgres"
	dbname   = "Forum"
)

// func createToken(c *gin.Context) {
// 	println("CREATING TOKEN", c)
// 	token := uuid.New().String()
// 	user := auth.User(c.Request)
// 	auth.Append(tokenStrategy, token, user)
// 	body := fmt.Sprintf("token: %s \n", token)
// 	c.String(http.StatusOK, body)
// }

// Middleware for authenticating requests
//

func parser(c *gin.Context, data interface{}, message string) {

	req, exists := c.Get("request")
	if !exists {
		// Handle error if the request is not found in context
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request not found in context"})
		return
	}

	if err := c.ShouldBindJSON(data); err != nil {
		// return this if we can't parse the data
		c.JSON(400, gin.H{"error": "Invalid request"})
	}
	c.JSON(400, gin.H{"Sucess": message})
}

func createSubPost(c *gin.Context) {
	var subPost struct {
		Creator_id  string `json:"creator_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	parser(c, &subPost, "Binded a subPost")

	Name := subPost.Name
	Description := subPost.Description
	Creator_Id := subPost.Creator_id
	println("CHECKING ID", Creator_Id)
	_, err := db.Exec("INSERT INTO subposts (name, description, creator_id) VALUES ($1, $2, $3)", Name, Description, Creator_Id)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to insert SubPost"})
		return
	}

	// Respond with success message or status
	c.JSON(200, gin.H{"message": "SubPost created successfully"})
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
	router.POST("/v1/auth/login", func(c *gin.Context) {
		Myauth.LoginHandler(c, db)
	})
	// router.POST("/v1/auth/refresh", Myauth.RefreshHandler)
	router.POST("v1/signup", Myauth.Signup)

	// router.GET("/v1/auth/token", middleware(createToken))

	// Start and run the server
	router.Run(":5000")

	fmt.Println("Server is running on http://localhost:5000")
	fmt.Println("Successfully connected!")
}
