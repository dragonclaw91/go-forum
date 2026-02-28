package repository

import (
	"strings"

	"github.com/gin-gonic/gin"

	"davidbrown/go/Go-Forum-App/cmd/internal/apperrs"
	"davidbrown/go/Go-Forum-App/cmd/internal/models"
)

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

	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

// pass the user name instead of a pointer because we are not mutating the data
func CheckUserExists(c *gin.Context, name string) error {
	var exists bool

	name = strings.TrimSpace(name)

	/* We are using the raw here, because if we used gorm,
	it would default to false, if nothing was found.
	Which might make it hard to disguingusih between
	broken database/sql string and a real false statement */

	err := db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)", name).
		Scan(&exists).Error

	if err != nil {
		return apperrs.Errgeneric
	}
	if exists {
		return apperrs.ErrUserNameTaken
	}

	return nil

}
