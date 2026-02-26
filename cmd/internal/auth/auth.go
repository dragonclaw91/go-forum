package Myauth

//TODO: we need to rig up the profile pic endpoint somewhere along the lines

//TODO: Look into GROM it might be a cleaner easier way to handle sql
import (
	"bytes"
	"database/sql"
	"davidbrown/go/Go-Forum-App/cmd/internal/apperrs"
	"davidbrown/go/Go-Forum-App/cmd/internal/models"
	"davidbrown/go/Go-Forum-App/cmd/internal/repository"
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

type Tokens struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// this will be used when creating a user largely
// type Credentials struct {
// 	Name     string
// 	Password string
// }

// var (
//
//	ErrUserNameTaken = errors.New("user name is not available")
//	ErrInvalidPass   = errors.New("username or password does not match")
//	ErrDatabaseDown  = errors.New("internal connection error")
//	ErrFailedToken   = errors.New("failed to generate token")
//	ErrBlankFields   = errors.New("Fields can not be blank")
//
// )
var strategy union.Union
var keeper jwt.SecretsKeeper

// JWT expiration times
var accessTokenExpiration = time.Minute * 15     // 15 minutes for access token
var refreshTokenExpiration = time.Hour * 24 * 30 // 30 days for refresh token

// TODO: Look into this to see if we can clean this up and just use Gin to read and parse
func Middleware(next gin.HandlerFunc) gin.HandlerFunc {

	//TODO: we moved the logic for the access token to a struct we need to get rid of this

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

func GetCredientials(c *gin.Context, db *sql.DB) (*models.Creds, error) {
	println("in the GetCredientials")
	var creds models.Creds

	if err := c.ShouldBindJSON(&creds); err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic in queryData:", r)
			}
		}()
	}
	if creds.Name == "" || creds.Password == "" {
		println("Passed ")
		fmt.Printf("%+v\n", creds)
		return nil, apperrs.ErrBlankFields
	} else {

		fmt.Printf("%+v\n", creds)
		return &creds, nil
	}

}

func GetUser(c *gin.Context, db *sql.DB, creds *models.Creds) (*repository.User, error) {

	result, err := repository.GetUser(creds.Name, c)

	if err != nil {
		return nil, apperrs.ErrInvalidPass
	}

	return result, nil

}

// LoginHandler generates both access and refresh tokens
func LoginHandler(c *gin.Context, db *sql.DB) {
	creds, err := GetCredientials(c, db)
	if err != nil {

		// If it's a blank field, we can be specific.
		if err.Error() == "Fields can not be blank" {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		} else {
			// If it's a technical JSON failure, stay generic.
			c.JSON(400, gin.H{"error": apperrs.Errgeneric.Error()})
			return
		}

	}

	result, err := GetUser(c, db, creds)

	if err != nil {
		// If it's a blank field, we can be specific.
		if err.Error() == "username or password does not match" {

			c.JSON(400, gin.H{"error": err.Error()})
			return
		} else {
			// If it's a technical JSON failure, stay generic.
			c.JSON(400, gin.H{"error": apperrs.Errgeneric.Error()})
			return
		}
	}

	// vaildUser is simply going to return true or false dependening if the correct creds were given

	err = validUser(result.Password, creds.Password)
	if err != nil {
		if err.Error() == "username or password does not match" {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		} else {
			// If it's a technical JSON failure, stay generic.
			c.JSON(400, gin.H{"error": apperrs.Errgeneric.Error()})
			return
		}
	}

	token, err := generateTokenSuite(result.Name)

	c.JSON(200, gin.H{
		"access_token": token.Access,
	})

}

// Validates the username and password (simple example)
func validUser(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password doesn't match
		fmt.Println("Incorrect password")
		return apperrs.ErrInvalidPass
	} else {
		// Password matches
		fmt.Println("Password match!")
		return nil

	}
}

// generateJWTToken creates a JWT token with the given expiration and secret
func generateJWTToken(username string, expiration time.Duration, secret []byte) (string, error) {

	//TODO: we need to add the access tokens to have role on this as well
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

// The Specialist Function
// TODO: this needs to be rigged up tested and implemented
func generateTokenSuite(username string) (*Tokens, error) {
	// 1. Create Access Token
	accessToken, err := generateJWTToken(username, accessTokenExpiration, jwtSecret)
	if err != nil {
		return nil, err // Send the error up the chain
	}

	// 2. Create Refresh Token
	refreshToken, err := generateJWTToken(username, refreshTokenExpiration, refreshSecret)
	if err != nil {
		return nil, err
	}

	// 3. Hand back the "Suite"
	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

// This needs to be returning something finsih this out
func setRefreshCookie(c *gin.Context, tokens *Tokens) {

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

	//TODO: This needs to be abstracted
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

	var creds *models.Creds

	//  parse the incoming JSON request and bind it to the creds struct.
	if err := c.ShouldBindJSON(&creds); err != nil {

		// return this if we can't parse the data or more likley the field was left blank
		c.JSON(400, gin.H{"error": apperrs.ErrBlankFields.Error()})
		return
	}
	// Checking for blank sapces to ensure some can't just bypass the form
	if strings.TrimSpace(creds.Name) == "" || strings.TrimSpace(creds.Password) == "" {
		c.JSON(400, gin.H{"error": apperrs.ErrBlankFields.Error()})
		return
	}

	result, err = repository.InsertUser(&models.Creds)

	// TODO:abstract this to the database file once created
	// var exists bool
	// // We use 'EXISTS' because it's faster than 'SELECT *'—it stops looking after it finds one match.
	// query := `SELECT EXISTS(SELECT 1 FROM users WHERE name=$1)`
	// creds.Name = strings.TrimSpace(creds.Name)

	// err := db.QueryRow(query, creds.Name).Scan(&exists)
	// if err != nil {
	// 	fmt.Println("Database check failed:", err)
	// 	c.JSON(500, gin.H{"error": "Internal server error"})
	// 	return
	// }

	// if exists {
	// 	fmt.Printf("Blocked: User %s already exists\n", creds.Name)
	// 	c.JSON(400, gin.H{"error": "Username is already taken"})
	// 	return
	// }
	// creds.Password = hashPassword(creds.Password)
	// if creds.Password == "Error Hashing" {
	// 	data := gin.H{
	// 		"message": "There was a problem",
	// 	}
	// 	c.JSON(http.StatusOK, data)
	// } else {
	// 	fmt.Printf("No Error: %+v\n", creds)

	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			log.Println("Panic in queryData:", r)
	// 		}
	// 	}()

	// 	// TODO:abstract this out to the datbase file once created

	// 	rows, err := db.Query(`INSERT INTO users (name, password) VALUES ($1, $2)`, creds.Name, creds.Password)

	// 	if err != nil {

	// 		println(err.Error())
	// 		log.Fatalf("Query error: %v", err)
	// 	}

	// 	defer rows.Close()
	// 	// Respond with success message or status
	// 	c.JSON(200, gin.H{"message": "User created successfully",
	// 		"User": creds.Name,
	// 	})
	// }
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
