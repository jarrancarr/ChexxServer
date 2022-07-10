package store

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	sslmode := os.Getenv("db_sslmode")

	//Build connection string
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		dbHost, username, dbName, sslmode, password)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	//db.Debug().AutoMigrate(&Account{}, &Comment{}, &Match{}, &Message{}, &Team{}, &Kingdom{}, &AI{}) //Database migration
	db.Debug().AutoMigrate(&User{}, &Match{}, &Message{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
