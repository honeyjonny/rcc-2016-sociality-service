package middleware

import "time"

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type PostForm struct {
	Content string `form:"content" json:"content" binding:"required"`
}

type UserDTO struct {
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
}

type PostDTO struct {
	Created time.Time `json:"created"`
	Content string    `json:"content"`
}
