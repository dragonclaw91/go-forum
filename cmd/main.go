package main

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/internal/auth"

	"fmt"

	// "github.com/shaj13/libcache"
	// _ "github.com/shaj13/libcache/fifo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

// const creditScoreMin = 500
// const creditScoreMax = 900

type Task struct {
	Id          string `json:"id"`
	Task        string `json:"task"`
	IsCompleted string `json:"iscompleted"`
}

type Result struct {
	Value []Task
	Err   error
}

// Global variables
var jwtSecret []byte     // Secret for access token
var refreshSecret []byte // Secret for refresh token

var db *sql.DB
var Rsecret string
var Asecret string

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
// func middleware(next gin.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.Println("Executing Auth Middleware")
// 		_, user, err := strategy.AuthenticateRequest(c.Request)
// 		if err != nil {
// 			c.JSON(401, gin.H{"error": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}
// 		log.Printf("User %s Authenticated\n", user.GetUserName())
// 		c.Set("user", user)
// 		next(c)
// 	}
// }

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
	Myauth.SetupGoGuardian()

	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/v1/auth/login", Myauth.LoginHandler)
	router.POST("/v1/auth/refresh", Myauth.RefreshHandler)
	router.POST("v1/signup", Myauth.Signup)
	// router.GET("/v1/auth/token", middleware(createToken))

	// Start and run the server
	router.Run(":5000")

	fmt.Println("Server is running on http://localhost:5000")
	fmt.Println("Successfully connected!")
}
