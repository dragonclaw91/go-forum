package models

type User struct {
	UserId   uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Password string `gorm:"not null" json:"-"` // This will never be sent to the frontend
	// ProfilePic   string `json:"profile_pic"`
	// Role         string `json:"role"`
}

type Creds struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
