package repository

import (
	"log"
)

type User struct {
	user_id  uint   `gorm:"primaryKey"`
	name     string `gorm:"uniqueIndex"`
	password string `json:"-"` // This will never be sent to the frontend
	// ProfilePic   string `json:"profile_pic"`
	// Role         string `json:"role"`
}

func GetUser(name string) (result string, error string) {
	var user User

	rows, err := db.Query(`SELECT user_id, password, name FROM "users" WHERE name= $1`, creds.Name)

	if err != nil {

		println(err.Error())
		log.Fatalf("Query error: %v", err)
	}

	defer rows.Close()
	if rows.Next() { // Iterate through the result set
		err := rows.Scan(&user.user_id, &user.password, &user.name) // Scan the result into the user struct
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
	}
	return &user, nil

}
