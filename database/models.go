package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type User struct {
	gorm.Model
	UserName string `sql:"unique_index"`
	Password string `sql:"index"`
}

type Post struct {
	gorm.Model
	UserID uint   `sql:"index"`
	Text   string `sql:"size:144"`
}

type Session struct {
	gorm.Model
	Cookie string `sql:"unique_index"`
	UserID uint   `sql:"index"`
}
