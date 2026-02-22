package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// const (
// 	host     = "localhost"
// 	port     = 5400
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "Forum"
// )

func InitDB(psqlInfo string) error {
	var err error
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// db, err = sql.Open("postgres", psqlInfo)
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	// if err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	return err
	// }
	return err
}
