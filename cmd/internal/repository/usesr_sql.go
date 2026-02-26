package repository

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"davidbrown/go/Go-Forum-App/cmd/internal/apperrs"
	"davidbrown/go/Go-Forum-App/cmd/internal/models"
)

type User struct {
	UserId   uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Password string `gorm:"not null" json:"-"` // This will never be sent to the frontend
	// ProfilePic   string `json:"profile_pic"`
	// Role         string `json:"role"`
}

func GetUser(name string, c *gin.Context) (*User, error) {
	var user User

	result := db.Select("user_id", "password", "name").Where("name = ?", name).Find(&user)

	if result.Error != nil {
		println(result.Error)
		return nil, apperrs.ErrInvalidPass
	}

	return &user, nil
}
func CheckUserExists(c *gin.Context, creds *models.Creds) error {
	var exists bool
	// We use 'EXISTS' because it's faster than 'SELECT *'—it stops looking after it finds one match.
	// query := `SELECT EXISTS(SELECT 1 FROM users WHERE name=$1)`

	creds.Name = strings.TrimSpace(creds.Name)

	err := db.Table("users").Select("Exists (SELECT 1)").Where("name = ?", creds.Name).Error

	if err != nil {
		println(err)
		return apperrs.Errgeneric
	}
	if exists {
		return apperrs.ErrUserNameTaken
	}
	if !exists {
		return nil
	}

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

		// TODO:abstract this out to the datbase file once created

		rows, err := db.Query(`INSERT INTO users (name, password) VALUES ($1, $2)`, creds.Name, creds.Password)

		if err != nil {

			println(err.Error())
			log.Fatalf("Query error: %v", err)
		}

		defer rows.Close()
		// Respond with success message or status
		c.JSON(200, gin.H{"message": "User created successfully",
			"User": creds.Name,
		})
	}
}
