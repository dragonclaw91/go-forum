package repository

import (
	"strings"

	"github.com/gin-gonic/gin"

	"davidbrown/go/Go-Forum-App/cmd/internal/apperrs"
	"davidbrown/go/Go-Forum-App/cmd/internal/models"
)

// type User struct {
// 	UserId   uint   `gorm:"primaryKey"`
// 	Name     string `gorm:"uniqueIndex"`
// 	Password string `gorm:"not null" json:"-"` // This will never be sent to the frontend
// 	// ProfilePic   string `json:"profile_pic"`
// 	// Role         string `json:"role"`
// }

//we are passing a pointer here because we are returning something from the database

func GetUser(name string, c *gin.Context) (*models.User, error) {
	var user models.User

	result := db.Select("user_id", "password", "name").Where("name = ?", name).Find(&user)

	if result.Error != nil {
		println(result.Error)
		return nil, apperrs.ErrInvalidPass
	}

	return &user, nil
}

/*
	the reason we are passing the strings rather than the pointer

is that we need to perserve the orgional password in case additonal checks need to be preformed
such as password length, strength etc.
*/
func InsertUser(c *gin.Context, name string, hashedPassword string) error {
	user := models.User{
		Name:     name,
		Password: hashedPassword,
	}
	// if creds.Password == "Error Hashing" {
	// 	data := gin.H{f
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

	// TODO:abstract this out to the datbase file once created

	// rows, err := db.Query(`INSERT INTO users (name, password) VALUES ($1, $2)`, creds.Name, creds.Password)

	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil

	// defer rows.Close()
	// // Respond with success message or status
	// c.JSON(200, gin.H{"message": "User created successfully",
	// 	"User": creds.Name,
	// })
}

// pass the user name instead of a pointer because we are not mutating the data
func CheckUserExists(c *gin.Context, name string) error {
	var exists bool
	// We use 'EXISTS' because it's faster than 'SELECT *'—it stops looking after it finds one match.
	// query := `SELECT EXISTS(SELECT 1 FROM users WHERE name=$1)`

	name = strings.TrimSpace(name)

	/* We are using the raw here, because if we used gorm,
	it would default to false, if nothing was found.
	Which might make it hard to disguingusih between
	broken database/sql string and a real false statement */

	err := db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", name).
		Scan(&exists).Error

	if err != nil {
		println(err)
		return apperrs.Errgeneric
	}
	if exists {
		return apperrs.ErrUserNameTaken
	}

	return nil

}
