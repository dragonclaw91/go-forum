package main

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	// "github.com/shaj13/libcache"
	// _ "github.com/shaj13/libcache/fifo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

// Global variables

var db *sql.DB

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "postgres"
	dbname   = "Forum"
)

// GLOBAL VARAIBLES
type QueryParams struct {
	Args     []interface{} // Parameters for the query (like where conditions)
	ScanArgs []interface{} // Variables to hold the scanned result
	Multi    []any
	Single   any
	Copy     any
}

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

	_, err := db.Exec(sql, data...)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to insert " + message})
		return
	}
	// Respond with success message or status
	c.JSON(200, gin.H{"message": message + " created successfully"})

}

func getHelper(c *gin.Context, sqlQuery string, isSingleRow bool, params QueryParams) (*sql.Row, *sql.Rows, error) {
	// Use reflection to check the type of Single
	v := reflect.ValueOf(params.Single)

	if isSingleRow {
		// If querying for a single row
		result := db.QueryRow(sqlQuery, params.Args...)

		err := result.Scan(params.ScanArgs...)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows found")
				return nil, nil, fmt.Errorf("no topic found with ID %s", params.Args...)
			}
			fmt.Println("Error scanning row:", err)
			return nil, nil, err
		}

		return result, nil, nil
	} else {
		// If querying for multiple rows
		result, err := db.Query(sqlQuery, params.Args...)
		params.Multi = []any{}

		for result.Next() {

			if err := result.Scan(params.ScanArgs...); err != nil {
				log.Fatal(err)
				return nil, nil, err
			}
			fmt.Println(result)

			// Check if Single is a pointer and if it's a *Thing Passed in
			if v.Kind() == reflect.Ptr {
				// Dereference the pointer and check if it points to a struct of type Topics
				if v.Elem().Kind() == reflect.Struct {
					// Dynamically append a copy of the struct to params.Multi
					// Make a copy of the struct (dereference the pointer and use reflect to make a copy)
					copyOfStruct := v.Elem().Interface()
					params.Multi = append(params.Multi, copyOfStruct)
				}
			}

		}
		c.JSON(http.StatusOK, params.Multi)
		if err != nil {
			log.Fatal(err)
		}
		return nil, result, nil
	}

	// return result
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

func getSubPost(c *gin.Context) {

	type Topics struct {
		ResultId    string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Creator_Id  string `json:"creator_id"`
		CreatedAt   string `json:"created_at"`
		Likes       string `json:"likes"`
		PostCount   string `json:"post_count"`
	}
	var topic Topics
	var topics []Topics

	sql := ""
	orderby := "ORDER BY created_at DESC"
	isSingleRow := false
	var subPost struct {
		Id      string `json:"id"`
		Search  string `json:"search"`
		OrderBy string `json:"order_by"`
	}

	parser(c, &subPost, "Failed to Bind a subPost")

	Id := subPost.Id
	Search := "%" + subPost.Search
	OrderBy := subPost.OrderBy

	params := QueryParams{
		Args: []interface{}{Id}, // SQL parameters (e.g., the ID)
		ScanArgs: []interface{}{
			&topic.ResultId,
			&topic.Name,
			&topic.Description,
			&topic.Creator_Id,
			&topic.CreatedAt,
			&topic.Likes},
	}

	if OrderBy != "" {
		switch {
		case OrderBy == "likes":
			orderby = `ORDER BY likes DESC`
		case OrderBy == "posts":
			orderby = `ORDER BY post_count DESC`
		}
	}

	switch {
	case Id != "":
		isSingleRow = true
		sql = `SELECT * FROM subposts WHERE id= $1 ` + orderby

		_, _, err := getHelper(c, sql, isSingleRow, params)
		if err == nil {
			c.JSON(http.StatusOK, topic)
		}

		if err != nil {
			println("FATAL ERROR", err.Error())
			log.Fatal(err)
		}

	case Search != "":
		params.Copy = topic
		params.Single = &topic
		params.Multi = []any{topics}
		params.Args = []interface{}{Search}
		params.ScanArgs = []interface{}{&topic.Description, &topic.CreatedAt, &topic.Likes, &topic.PostCount}
		sql = `SELECT s.name,s.created_at,s.likes, COUNT(p.sub_post_id) AS post_count
			FROM subposts s
			LEFT JOIN posts p ON p.sub_post_id = s.id
			WHERE s.name ILIKE $1
			GROUP BY s.id ` + orderby
		_, _, err := getHelper(c, sql, isSingleRow, params)
		if err != nil {
			println("FATAL ERROR", err.Error())
			log.Fatal(err)
		}
	default:
		params.Copy = topic
		params.Single = &topic
		params.Multi = []any{topics}
		params.Args = []interface{}{Search}
		params.ScanArgs = []interface{}{&topic.Description, &topic.CreatedAt, &topic.Likes, &topic.PostCount}
		sql = `SELECT s.name,s.created_at,s.likes, COUNT(p.sub_post_id) AS post_count
			FROM subposts s
			LEFT JOIN posts p ON p.sub_post_id = s.id
			GROUP BY s.id ` + orderby
		getHelper(c, sql, isSingleRow, params)
	}
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

func createVote(c *gin.Context) {
	var Vote struct {
		User_id  string `json:"user_id"`
		PostId   string `json:"post_id"`
		ReplyId  string `json:"reply_id"`
		VoteType string `json:"vote_type"`
	}
	// bind the data to the struct
	parser(c, &Vote, "Failed to Bind a post")

	User_Id := Vote.User_id
	PostId := Vote.PostId
	RepyId := Vote.ReplyId
	VoteType := Vote.VoteType

	sql := "INSERT INTO votes (user_id, post_id, reply_id, vote_type) VALUES ($1, $2, $3, $4)"

	// update the database
	postHelper(c, "Vote", sql, User_Id, PostId, RepyId, VoteType)

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

	Myauth.SetupGoGuardian()

	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/v1/subpost/create", Myauth.Middleware(createSubPost))
	router.POST("/v1/post/create", Myauth.Middleware(createPost))
	router.POST("/v1/replies/create", Myauth.Middleware(createReply))
	router.POST("/v1/vote/create", Myauth.Middleware(createVote))

	router.GET("/v1/subpost", Myauth.Middleware(getSubPost))

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
