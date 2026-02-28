package main

// TODO: abstract all of the helper functions to their own db file

import (
	"database/sql"
	Myauth "davidbrown/go/Go-Forum-App/cmd/internal/auth"
	"davidbrown/go/Go-Forum-App/cmd/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	// "github.com/shaj13/libcache"
	// _ "github.com/shaj13/libcache/fifo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"

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

	result, err := db.Exec(sql, data...)
	if err != nil {
		println("MADE IT")
		log.Fatal(err)
		c.JSON(500, gin.H{"error": "Failed to insert " + message})
		return
	}
	// Respond with success message or status
	c.JSON(200, gin.H{"message": message + " created successfully"})

	c.Set("deleteResult", result)

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
		if err != nil {
			println("ERROR", err.Error(), result)
		}

		// fmt.Println(result)
		params.Multi = []any{}

		for result.Next() {

			if err := result.Scan(params.ScanArgs...); err != nil {
				log.Fatal(err)
				println(err.Error())
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
			println(err.Error())
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

func putSubPost(c *gin.Context) {
	// place to store the elements of key value pairs of things we want to update
	type Field struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	// storing an array of things we want to update plus its id
	type subPost struct {
		Id     string  `json:"id"`
		Fields []Field `json:"fields"`
	}

	var subpost subPost
	var values []any

	// binding all data to a supost type
	parser(c, &subpost, "Failed to Bind a subPost")

	// this should never change so we hard code it
	sql := "UPDATE subposts SET "

	// Iterate over the subpost to parse it to build the final sql string and the new values of things we want to change
	for i := 0; i < len(subpost.Fields); i++ {
		// append the name of the field we want to update plus a placeholder for postgres
		sql = sql + subpost.Fields[i].Key + " = $" + strconv.Itoa(i+1)
		// if we are not at the end of our loop add a comma to prepare for the next field
		if i+1 != len(subpost.Fields) {
			sql = sql + ", "
			// if we are at the end of our loop add the where clause and a final placeholder
		} else {
			sql = sql + " WHERE Id = $" + strconv.Itoa(i+2)
		}
		// shoveing the new values we want to use to update to later be passed as args
		values = append(values, subpost.Fields[i].Value)
	}
	// finnaly we add the id of the thing we want to update
	values = append(values, subpost.Id)
	fmt.Println(values...)
	// update the database
	postHelper(c, "SubPost", sql, values...)

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

	// links back to subPost not topics
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

func getReplies(c *gin.Context) {

	type Replies struct {
		ResultId      string  `json:"id"`
		ParentReplyId *string `json:"parent_reply_id,omitempty"`
		Reply         string  `json:"reply"`
		CreatedAt     string  `json:"created_at"`
		UserId        string  `json:"name"`
		Score         string  `json:"score"`
	}
	var reply Replies
	var replies []Replies

	sql := ""
	orderby := "ORDER BY c.created_at"
	isSingleRow := false
	var replyArgs struct {
		Id      string `json:"id"`
		Search  string `json:"search"`
		OrderBy string `json:"order_by"`
		ReplyId string `json:"reply_id"`
	}
	println("PROBABLY MADE IT HERE")
	parser(c, &replyArgs, "Failed to Bind a Reply")

	ReplyId := replyArgs.ReplyId
	Id := replyArgs.Id
	Search := "%" + replyArgs.Search
	OrderBy := replyArgs.OrderBy

	params := QueryParams{
		Args: []interface{}{Id}, // SQL parameters (e.g., the ID)
		ScanArgs: []interface{}{
			&reply.ResultId,
			&reply.ParentReplyId,
			&reply.Reply,
			&reply.CreatedAt,
			&reply.UserId,
			&reply.Score,
		},
	}

	if OrderBy != "" {
		switch {

		case OrderBy == "score":
			orderby = `ORDER BY score DESC`
		}
	}

	switch {
	case ReplyId != "":
		params.Single = &reply
		params.Multi = []any{replies}
		params.Args = []interface{}{Id, ReplyId}
		sql = `SELECT 
		c.id, parent_reply_id, reply, c.created_at, u.name,
			COALESCE(SUM(CASE WHEN v.vote_type = 'upvote' THEN 1 ELSE 0 END), 0) - 
			COALESCE(SUM(CASE WHEN v.vote_type = 'downvote' THEN 1 ELSE 0 END), 0) AS score
		FROM replies c
		JOIN users u ON c.user_id = u.user_id
		LEFT JOIN votes v ON c.id = v.reply_id
		WHERE c.parent_reply_id = $2 AND c.subpost_id = $1
		GROUP BY c.id, u.name ` + orderby

		_, _, err := getHelper(c, sql, isSingleRow, params)

		if err != nil {
			println("FATAL ERROR", err.Error())
			log.Fatal(err)
		}

	case Search != "":
		params.Copy = reply
		params.Single = &reply
		params.Multi = []any{replies}
		params.Args = []interface{}{Id, Search}
		sql = `SELECT 
		c.id, parent_reply_id, reply, c.created_at, u.name,
			COALESCE(SUM(CASE WHEN v.vote_type = 'upvote' THEN 1 ELSE 0 END), 0) - 
			COALESCE(SUM(CASE WHEN v.vote_type = 'downvote' THEN 1 ELSE 0 END), 0) AS score
		FROM replies c
		JOIN users u ON c.user_id = u.user_id
		LEFT JOIN votes v ON c.id = v.reply_id
		WHERE c.parent_reply_id IS NULL AND c.subpost_id = $1 AND reply ILIKE $2
		GROUP BY c.id, u.name ` + orderby
		_, _, err := getHelper(c, sql, isSingleRow, params)
		if err != nil {
			println("FATAL ERROR", err.Error())
			log.Fatal(err)
		}
	default:
		params.Single = &reply
		params.Multi = []any{replies}
		params.Args = []interface{}{Id}
		sql = `SELECT c.id, parent_reply_id,  reply, c.created_at, 
		COALESCE(SUM(CASE WHEN vote_type = 'upvote' THEN 1 ELSE 0 END), 0) - 
		COALESCE(SUM(CASE WHEN vote_type = 'downvote' THEN 1 ELSE 0 END), 0) AS score
 		FROM replies c
 		LEFT JOIN votes v ON c.id = v.reply_id
 		WHERE c.parent_reply_id IS NULL AND c.subpost_id = $1
 		GROUP BY c.id ` + orderby

		_, _, err := getHelper(c, sql, isSingleRow, params)

		if err != nil {
			println("FATAL ERROR", err.Error())
			log.Fatal(err)
		}
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

func putReply(c *gin.Context) {
	// place to store the new reply
	var Reply struct {
		Id      string `json:"id"`
		Reply   string `json:"reply"`
		Flagged bool   `json:"flagged"`
	}

	// binding all data to a supost type
	parser(c, &Reply, "Failed to Bind a reply")

	sql := "UPDATE replies SET reply = $1 WHERE id = $2 "
	switch {
	case Reply.Flagged:
		sql = `UPDATE replies
		SET flagged = NOT flagged
		WHERE id = $1;`
		postHelper(c, "Reply", sql, Reply.Id)

	default:
		sql = "UPDATE replies SET reply = $1 WHERE id = $2 "
		postHelper(c, "Reply", sql, Reply.Reply, Reply.Id)
	}

	// update the database

}

func DeleteHelper(c *gin.Context) {
	// place to store the new reply
	var Delete struct {
		Id         string `json:"id"`
		DeleteFrom string `json:"delete_from"`
	}

	// binding all data to a supost type
	parser(c, &Delete, "Failed to Bind a subpost")

	// this should never change so we hard code it
	sql := "DELETE FROM " + Delete.DeleteFrom + " WHERE id = $1"
	println(sql)
	// update the database
	postHelper(c, Delete.DeleteFrom+" delete", sql, Delete.Id)

}

func putHelper(c *gin.Context) {
	var Update struct {
		Id      string `json:"id"`
		Change  string `json:"change"`
		NewRole string `json:"new_role"`
		Picture string `json:"picture"`
	}

	// binding all data to a supost type
	parser(c, &Update, "Failed to Bind a subpost")

	// this should never change so we hard code it
	sql := ""

	switch {

	case Update.Change == "new_admin" || Update.Change == "new_moderator":
		sql = `UPDATE users SET pending = false `
		sql = sql + ` , ` + Update.NewRole + ` = true  WHERE user_id = $1;`
		println("SQL", sql, Update.Id)
		postHelper(c, "updated user", sql, Update.Id)

	case Update.Change == "picture":
		sql = `UPDATE users SET picture = $2 WHERE user_id = $1`
		if Update.Picture == "none" {
			postHelper(c, "updated user", sql, Update.Id, null.IntFromPtr(nil))
		} else {
			postHelper(c, "updated user", sql, Update.Id, Update.Picture)
		}
	}
}

// func init() {

// 	var err error
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)
// 	db, err = sql.Open("postgres", psqlInfo)

// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	err := repository.InitDB(psqlInfo)

	if err != nil {
		// We use log.Fatal because if the DB is down, the app is useless
		log.Fatalf("Could not connect to database: %v", err)
	}

	log.Println("Database connection established!")
	Myauth.SetupGoGuardian()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		// 1. MUST match your Angular URL exactly
		AllowOrigins: []string{"http://localhost:4200"},

		// 2. Allow the browser to see these specific headers
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// 3. THIS IS THE BIG ONE - This makes the headers "visible"
		AllowCredentials: true,

		// 4. How long the browser should trust this "Yes" (prevents constant OPTIONS checks)
		MaxAge: 12 * time.Hour,
	}))

	router.POST("/v1/subpost/create", Myauth.Middleware(createSubPost))
	router.POST("/v1/replies/create", Myauth.Middleware(createReply))
	router.POST("/v1/vote/create", Myauth.Middleware(createVote))
	router.POST("/v1/auth/login", func(c *gin.Context) {
		println("SHOULD BE THE FIRST STEP")
		c.Set("request", c.Request)
		Myauth.LoginHandler(c, db)

	})
	router.POST("/v1/auth/refresh", Myauth.RefreshHandler)
	router.POST("v1/auth/signup", func(c *gin.Context) {
		c.Set("request", c.Request)
		Myauth.Signup(c, db)
	})

	router.GET("/v1/subpost", Myauth.Middleware(getSubPost))
	router.GET("/v1/replies", Myauth.Middleware(getReplies))

	router.PUT("/v1/subpost/update", Myauth.Middleware(putSubPost))
	router.PUT("/v1/replies/update", Myauth.Middleware(putReply))
	router.PUT("/v1/user/update", Myauth.Middleware(putHelper))

	router.DELETE("/v1/delete", Myauth.Middleware(DeleteHelper))

	// router.GET("/v1/auth/token", middleware(createToken))

	// Start and run the server
	router.Run(":5000")

	fmt.Println("Server is running on http://localhost:5000")
	fmt.Println("Successfully connected!")
}
