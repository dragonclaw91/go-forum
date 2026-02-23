package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	UserId   uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Password string `gorm:"not null" json:"-"` // This will never be sent to the frontend
	// ProfilePic   string `json:"profile_pic"`
	// Role         string `json:"role"`
}

func GetUser(name string, c *gin.Context) *gorm.DB {
	var user User
	// name = "test"
	println("MADE IT TOTHE FUNCTION")
	fmt.Printf("DEBUG: Target struct state before GORM: %+v\n", user)
	// var user User
	// ctx := c.Request.Context()
	// product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx)
	// user, err := gorm.G[User](db).Where("name = ?",name).First(ctx)
	result := db.First(&user, "name = ?", name)
	// rows, err := db.Query(`SELECT user_id, password, name FROM "users" WHERE name= $1`, name)

	if result.Error != nil {

		println(result.Error)
		return result
		// log.Fatalf("Query error: %v", result.Error)
	}

	// defer rows.Close()
	// if rows.Next() { // Iterate through the result set
	// 	err := rows.Scan(&user.user_id, &user.password, &user.name) // Scan the result into the user struct
	// 	if err != nil {
	// 		log.Fatalf("Error scanning row: %v", err)
	// 	}
	// }
	// return &user, nil
	// return result
	fmt.Printf("Fetched User: %+v\n", user)
	return result
}
