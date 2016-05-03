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

	db, err := gorm.Open(config.Dialect, config.ConnectionString) //"postgres", "user=hj dbname=undesire password=hJ073 sslmode=disable")

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
}

func set_constraints(db *gorm.DB) {
	db.Model(&Session{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}
