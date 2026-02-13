package Myauth

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	dgjwt "github.com/golang-jwt/jwt"

	// "github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"golang.org/x/crypto/bcrypt"
)

// Global variables
var jwtSecretString = os.Getenv("ACCESS_TOKEN_SECRET")
var jwtSecret = []byte(jwtSecretString) // Secret for access token

var refreshSecretString = os.Getenv("REFRESH_TOKEN_SECRET")
var refreshSecret = []byte(refreshSecretString) // Secret for refresh token

// var db *sql.DB

// var Rsecret string
// var Asecret string
// var tokenStrategy auth.Strategy
// var cacheObj libcache.Cache
var strategy union.Union
var keeper jwt.SecretsKeeper

// JWT expiration times
var accessTokenExpiration = time.Minute * 15     // 15 minutes for access token
var refreshTokenExpiration = time.Hour * 24 * 30 // 30 days for refresh token

func Middleware(next gin.HandlerFunc) gin.HandlerFunc {

	// Struct to match the expected JSON data
	type RequestData struct {
		AccessToken string `json:"access_token"`
	}

	return func(c *gin.Context) {

		var requestData RequestData

		// Read the request body into a buffer
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read request body"})
			c.Abort()
			return
		}

		// Reassign the original request body to the new body (for further use)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Store the body in context for later use
		c.Set("rawBody", body)

		// Parse JSON body into the requestData struct
		if err := c.ShouldBindJSON(&requestData); err != nil {
			// If there’s an error unmarshalling, respond with an error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		// Print  access_token
		fmt.Println("Access TokenZ:", requestData.AccessToken)
		println("IN THE MIDDLEWARE", c.DefaultPostForm("access_token", ""))

		_, err = ValidateJWT(requestData.AccessToken, jwtSecret)
		if err != nil {
			println("ERROR", err)
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("Executing Auth Middleware")
		next(c)
	}
}

func SetupGoGuardian() {
	keeper = jwt.StaticSecret{
		Secret:    jwtSecret,
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 15)
	jwtStrategy := jwt.New(cache, keeper)
	strategy = union.New(jwtStrategy)
}

// LoginHandler generates both access and refresh tokens
func LoginHandler(c *gin.Context, db *sql.DB) {
	println("in the LoginHandler")
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

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic in queryData:", r)
		}
	}()

	rows, err := db.Query(`SELECT password FROM "users" WHERE name= $1`, name)

	if err != nil {

		println(err.Error())
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

		accessToken, err := generateJWTToken(creds.Name, accessTokenExpiration, jwtSecret)
		if err != nil {
			// if we cant generate a access token return this
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			// return
		}

		//Generate refresh token
		// Generate secure random JWT secrets for both access and refresh tokens
		refreshToken, err := generateJWTToken(creds.Name, refreshTokenExpiration, refreshSecret)
		if err != nil {
			// if we cant generate a refresh token return this
			c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
			// return
		}
		println("Setting g the cookie")
		c.SetCookie(
			"refresh_token",
			refreshToken,
			3600*24*30, // 30 days
			"/",
			"",
			false, // Set to true when you have HTTPS/SSL
			true,  // The "Pro" flag: HttpOnly
		)
		println("Cookie set")

		// Send both tokens to the client
		c.JSON(200, gin.H{
			"access_token": accessToken,
		})
	} else {
		// failing to verfiy username and password return this
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	}
}

// Validates the username and password (simple example)
func validUser(hashedPassword, password string) bool {
	fmt.Printf("DEBUG: Hashed from DB: %q\n", hashedPassword)
	fmt.Printf("DEBUG: Plain from User: %q\n", password)

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
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// REFRESH FUNCTIONS

func RefreshHandler(c *gin.Context) {
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

	// Validate the refresh token with the refresh secret
	claims, err := ValidateJWT(req.RefreshToken, refreshSecret)
	// println("ERRR", err.Error())
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token again"})
		return
	}
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

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		3600*24*30, // 30 days
		"/",
		"",
		false, // Set to true when you have HTTPS/SSL
		true,  // The "Pro" flag: HttpOnly
	)

	// Return the new access token
	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}

// validateJWT validates the JWT token (either access or refresh)
func ValidateJWT(tokenString string, secret []byte) (dgjwt.MapClaims, error) {
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

// SIGNUP FUNCTIONS

func Signup(c *gin.Context, db *sql.DB) {

	var creds struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	//  parse the incoming JSON request and bind it to the creds struct.
	if err := c.ShouldBindJSON(&creds); err != nil {

		// return this if we can't parse the data or more likley the field was left blank
		c.JSON(400, gin.H{"error": "All fields are required!"})
		return
	}
	// Checking for blank sapces to ensure some can't just bypass the form
	if strings.TrimSpace(creds.Name) == "" || strings.TrimSpace(creds.Password) == "" {
		c.JSON(400, gin.H{"error": "Username and password cannot be empty"})
		return
	}
	var exists bool
	// We use 'EXISTS' because it's faster than 'SELECT *'—it stops looking after it finds one match.
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE name=$1)`
	creds.Name = strings.TrimSpace(creds.Name)

	err := db.QueryRow(query, creds.Name).Scan(&exists)
	if err != nil {
		fmt.Println("Database check failed:", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if exists {
		fmt.Printf("Blocked: User %s already exists\n", creds.Name)
		c.JSON(400, gin.H{"error": "Username is already taken"})
		return
	}
	creds.Password = hashPassword(creds.Password)
	if creds.Password == "Error Hashing" {
		data := gin.H{
			"message": "There was a problem",
		}
		c.JSON(http.StatusOK, data)
	} else {
		fmt.Printf("No Error: %+v\n", creds)

		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic in queryData:", r)
			}
		}()

		rows, err := db.Query(`INSERT INTO users (name, password) VALUES ($1, $2)`, creds.Name, creds.Password)

		if err != nil {

			println(err.Error())
			log.Fatalf("Query error: %v", err)
		}

		defer rows.Close()
		// Respond with success message or status
		c.JSON(200, gin.H{"message": "User created successfully"})
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
