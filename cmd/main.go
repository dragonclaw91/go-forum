package main

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/internal/auth"
	"encoding/json"
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

		c.Set("request", c.Request)
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
