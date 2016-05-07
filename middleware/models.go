package middleware

import "time"

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type MessageForm struct {
	Uid     string `form:"uid" json:"uid" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

type PostForm struct {
	Content string `form:"content" json:"content" binding:"required"`
}

type UserDTO struct {
	Uid      uint      `json:"uid"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
}

type PostDTO struct {
	Created time.Time `json:"created"`
	Content string    `json:"content"`
}

type MessageDTO struct {
	Created time.Time `json:"created"`
	Message string    `json:"message"`
}

type NodeDTO struct {
	Id    uint   `json:"id"`
	Label string `json:"label"`
}

type LinkDTO struct {
	From uint `json:"from"`
	To   uint `json:"to"`
}
