package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	dgjwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"golang.org/x/crypto/bcrypt"
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
var tokenStrategy auth.Strategy
var cacheObj libcache.Cache
var strategy union.Union
var keeper jwt.SecretsKeeper

// JWT expiration times
var accessTokenExpiration = time.Minute * 15     // 15 minutes for access token
var refreshTokenExpiration = time.Hour * 24 * 30 // 30 days for refresh token

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "postgres"
	dbname   = "Forum"
)

func createToken(c *gin.Context) {
	println("CREATING TOKEN", c)
	token := uuid.New().String()
	user := auth.User(c.Request)
	auth.Append(tokenStrategy, token, user)
	body := fmt.Sprintf("token: %s \n", token)
	c.String(http.StatusOK, body)
}

// Generate secure random JWT secrets for both access and refresh tokens
func generateRandomJWTSecret() error {
	// Generate a secret for the access token (JWT secret)
	accessTokenSecret := make([]byte, 32) // 256 bits = 32 bytes
	_, err := rand.Read(accessTokenSecret)
	if err != nil {
		return fmt.Errorf("failed to generate access token secret: %w", err)
	}
	jwtSecret = accessTokenSecret

	// Generate a separate secret for the refresh token
	refreshTokenSecret := make([]byte, 32) // 256 bits = 32 bytes
	_, err = rand.Read(refreshTokenSecret)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token secret: %w", err)
	}
	refreshSecret = refreshTokenSecret

	// Log the secrets (for debugging purposes, but remember to never log in production)
	// log.Printf("Generated JWT Secret (Base64): %s", base64.StdEncoding.EncodeToString(jwtSecret))
	// log.Printf("Generated Refresh Secret (Base64): %s", base64.StdEncoding.EncodeToString(refreshTokenSecret))

	return nil
}

// LoginHandler generates both access and refresh tokens
func loginHandler(c *gin.Context) {
	var creds struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	//  parse the incoming JSON request and bind it to the creds struct.
	if err := c.ShouldBindJSON(&creds); err != nil {
		// return this if we can't parse the data
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	name := creds.Name
	password := creds.Password
	rows, err := db.Query(`SELECT password FROM "users" WHERE name= $1`, name)

	// err := db.Query()
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	defer rows.Close()
	if rows.Next() { // Iterate through the result set
		err := rows.Scan(&password) // Scan the result into the password variable
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
	}
	println("CHECKING RESULT", name, password)

	// call validUser to verify username and password
	// vaildUser is simply going to return true or false dependening if the correct creds were given
	if validUser(password, creds.Password) {
		/*if vaild username and password  generate access and refresh tokens
		handled by the generate JWTToken function */
		// Generate secure random JWT secrets for both access and refresh tokens
		if err := generateRandomJWTSecret(); err != nil {
			log.Fatalf("Error generating JWT secret: %v", err)
		}
		accessToken, err := generateJWTToken(creds.Name, accessTokenExpiration, jwtSecret)
		if err != nil {
			// if we cant generate a access token return this
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		//Generate refresh token
		// Generate secure random JWT secrets for both access and refresh tokens
		// if err := generateRandomJWTSecret(); err != nil {
		// 	log.Fatalf("Error generating JWT secret: %v", err)
		// }
		refreshToken, err := generateJWTToken(creds.Name, refreshTokenExpiration, refreshSecret)
		if err != nil {
			// if we cant generate a refresh token return this
			c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		// Send both tokens to the client
		c.JSON(200, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	} else {
		// failing to verfiy username and password return this
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	}
}

// Validates the username and password (simple example)
func validUser(hashedPassword, password string) bool {
	println("VALIDATING")

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password doesn't match
		fmt.Println("Incorrect password")
		return false
	} else {
		// Password matches
		fmt.Println("Password match!")
		return true
	}
	// if err := c.ShouldBindJSON(&login); err != nil {
	// 	// if we can't parse the json return this
	// 	log.Printf("Error binding JSON: %v", err)
	// 	c.JSON(400, gin.H{"error": "Invalid request"})
	// 	return false
	// }
	// // err := db.Query()
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// println("CHECKING RESULT", login)
	// TODO rig up an actual query to the database
	// using dummy data looking for these valuse in the request
	// return false
}

// generateJWTToken creates a JWT token with the given expiration and secret
func generateJWTToken(username string, expiration time.Duration, secret []byte) (string, error) {

	// setting the expriation of the token and the username associated with it
	claims := dgjwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(expiration).Unix(), // Set expiration
	}
	// create the token
	token := dgjwt.NewWithClaims(dgjwt.SigningMethodHS256, claims)
	// Sign the token with the passed secret
	// Generate secure random JWT secrets for both access and refresh tokens
	// if err := generateRandomJWTSecret(); err != nil {
	// 	log.Fatalf("Error generating JWT secret: %v", err)
	// }
	log.Printf("CREATE secret: %s", base64.StdEncoding.EncodeToString(secret))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func refreshHandler(c *gin.Context) {
	// var refreshToken string
	type RefreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	var req RefreshRequest

	// binding the refresh token
	if err := c.ShouldBindJSON(&req); err != nil {
		// if we can't parse the json return this
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	// refreshToken = req.RefreshToken
	println("TOKEN", req.RefreshToken)
	log.Printf("Base64 encoded secret: %s", base64.StdEncoding.EncodeToString(refreshSecret))
	// log.Printf("JWT 1Secret (Base64): %s", base64.StdEncoding.EncodeToString(refreshSecret))

	// Validate the refresh token with the refresh secret
	claims, err := validateJWT(req.RefreshToken, refreshSecret)
	// println("ERRR", err.Error())
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token again"})
		return
	}
	println("MOVING ON")
	// Generate a new access token if the refresh token is valid
	accessToken, err := generateJWTToken(claims["sub"].(string), accessTokenExpiration, jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate new access token"})
		return
	}

	// Generate a new refresh token
	newRefreshToken, err := generateJWTToken(claims["sub"].(string), refreshTokenExpiration, refreshSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	// Return the new access token
	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

// validateJWT validates the JWT token (either access or refresh)
func validateJWT(tokenString string, secret []byte) (dgjwt.MapClaims, error) {
	// return the secret to vierify signature
	log.Printf("JWT Secret (Base64): %s", base64.StdEncoding.EncodeToString(secret))
	log.Printf("JWT Token: %s", tokenString)

	token, err := dgjwt.Parse(tokenString, func(token *dgjwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		println("ERROR", err.Error())
		log.Printf("Invalid JWT token: %v", token)
		// if the signature can't be verifed return this
		return nil, fmt.Errorf("invalid token")
	}
	// create the claims
	claims, ok := token.Claims.(dgjwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func setupGoGuardian() {
	keeper = jwt.StaticSecret{
		Secret:    jwtSecret,
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)
	jwtStrategy := jwt.New(cache, keeper)
	strategy = union.New(jwtStrategy)
}

// Middleware for authenticating requests
func middleware(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Executing Auth Middleware")
		_, user, err := strategy.AuthenticateRequest(c.Request)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		log.Printf("User %s Authenticated\n", user.GetUserName())
		c.Set("user", user)
		next(c)
	}
}

func hashPassword(password string) string {
	println("Pre Hash", password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err == nil {
		return string(bytes)
	} else {
		println("ERROR", err)
		return "Error Hashing"
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

func taskGet() Result {

	rows, err := db.Query(`SELECT name,password FROM "users"`)
	// err := db.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tasks []Task
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Task, &task.IsCompleted); err != nil {
			fmt.Println("<<<<<<< In the Get Function >>>>>>>>>", err)
			return Result{tasks, err}
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return Result{tasks, err}
	}

	return Result{tasks, nil}

}

func taskPostPut(c *gin.Context, sql string, arg string, ids string) {
	var err error
	// Handle the POST request...
	// stmt, err :=
	// if err != nil {
	// 	log.Println(err)
	// 	return false
	//    }
	//    defer stmt.Close()
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	// fmt.Println(body)
	println("CHECKING THE BODY", string(body))
	// if err != nil {
	// 	http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 	return
	// }
	defer c.Request.Body.Close()

	// Parse the JSON body
	// var task Task
	// var user Login

	// if ids == "" {
	// 	fmt.Println("New DELETE:", string(body))
	// 	if err := json.Unmarshal(body, &user); err != nil {
	// 		fmt.Println("Error unmarshalling JSON:", err)
	// 		return
	// 	}
	// }
	sqlStatement := sql

	id := ""
	if arg == "POST" {
		// err = db.QueryRow(`SELECT name FROM users WHERE name=$1`, user.Name).Scan(&id)
		if err == nil {
			data := gin.H{
				"message": "There was a problem",
			}
			c.JSON(http.StatusOK, data)
		} else {
			// user.Password = hashPassword(user.Password)
			// if user.Password == "Error Hashing" {
			// 	data := gin.H{
			// 		"message": "There was a problem",
			// 	}
			// 	c.JSON(http.StatusOK, data)
			// } else {
			// 	err = db.QueryRow(sqlStatement, user.Name, user.Password).Scan(&id)
			// }
		}

	}
	// if arg == "PUT" {
	// 	// fmt.Println("New DELETE:", task)
	// 	err = db.QueryRow(sqlStatement, task.Id).Scan(&id)
	// }
	if arg == "DELETE" {
		// fmt.Println("New DELETE:", ids)
		err = db.QueryRow(sqlStatement, ids).Scan(&id)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", arg)

}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("<<<<<<< In the handler >>>>>>>>>")
// 	switch r.Method {
// 	case http.MethodGet:
// 		taskGet()
// 	case http.MethodPost:
// 		taskPostPut(w, r, `
// 		INSERT INTO "List" ("task")
// 		VALUES ($1) RETURNING id`, "POST")
// 		fmt.Fprintln(w, "You made a POST request!")
// 	case http.MethodPut:
// 		taskPostPut(w, r, `UPDATE "List" SET "isComplete"=NOT "isComplete","completed_at"=NOW() WHERE "id"=$1 RETURNING id`, "PUT")
// 		fmt.Fprintln(w, "You made a PUT request!")
// 	case http.MethodDelete:
// 		fmt.Fprintln(w, "You made a DELETE request!")
// 		taskPostPut(w, r, ` DELETE FROM "List" WHERE "id"=$1 RETURNING id`, "PUT")
// 	default:
// 		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
// 	}
// }

func signup(c *gin.Context) {

	var creds struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	//  parse the incoming JSON request and bind it to the creds struct.
	if err := c.ShouldBindJSON(&creds); err != nil {
		// return this if we can't parse the data
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	creds.Password = hashPassword(creds.Password)
	if creds.Password == "Error Hashing" {
		data := gin.H{
			"message": "There was a problem",
		}
		c.JSON(http.StatusOK, data)
	} else {
		_, err := db.Exec("INSERT INTO users (name, password) VALUES ($1, $2)", creds.Name, creds.Password)
		if err != nil {
			log.Fatal(err)
			c.JSON(500, gin.H{"error": "Failed to insert user"})
			return
		}

		// Respond with success message or status
		c.JSON(200, gin.H{"message": "User created successfully"})
	}
}

func main() {

	setupGoGuardian()
	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/v1/auth/login", loginHandler)
	router.POST("/v1/auth/refresh", refreshHandler)
	router.POST("v1/signup", signup)
	router.GET("/v1/auth/token", middleware(createToken))

	// Start and run the server
	router.Run(":5000")
	// fmt.Println("New record ID is:", id)
	// http.HandleFunc("/", handler)

	fmt.Println("Server is running on http://localhost:5000")
	// http.ListenAndServe(":8080", nil)
	fmt.Println("Successfully connected!")
}
