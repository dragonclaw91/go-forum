package repository

import (
	"github.com/gin-gonic/gin"
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
		return nil, db.Error
	}

	return &user, nil
}
