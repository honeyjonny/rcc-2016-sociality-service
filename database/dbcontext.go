package database

import (
	_ "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DbConfig struct {
	Dialect          string
	ConnectionString string
}

func (config *DbConfig) CreateConnection() *gorm.DB {

	db, err := gorm.Open(config.Dialect, config.ConnectionString)

	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	migrate(db)
	set_constraints(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Follower{})
	db.AutoMigrate(&Message{})
}

func set_constraints(db *gorm.DB) {
	db.Model(&Session{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&Follower{}).AddForeignKey("object_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&Follower{}).AddForeignKey("subject_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&Message{}).AddForeignKey("object_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&Message{}).AddForeignKey("subject_id", "users(id)", "CASCADE", "CASCADE")
}
